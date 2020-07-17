package convert

import (
	"io"
	"sort"

	gotoml "github.com/pelletier/go-toml"

	"github.com/sclevine/yj/order"
)

type TOML struct {
	Indent bool
}

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) error {
	enc := gotoml.NewEncoder(&trimWriter{w: w})
	enc.Order(gotoml.OrderPreserve)
	if !t.Indent {
		enc.Indentation("")
	}
	return enc.Encode(jsonToTOML(in))
}

type trimWriter struct {
	w    io.Writer
	done bool
}

func (w *trimWriter) Write(p []byte) (n int, err error) {
	trimmed := false
	if !w.done && len(p) > 0 && p[0] == '\n' {
		p = p[1:]
		trimmed = true
	}
	n, err = w.w.Write(p)
	if (trimmed && err == nil) || n > 0 {
		w.done = true
		if trimmed {
			n++
		}
	}
	return n, err
}

func jsonToTOML(in interface{}) interface{} {
	switch in := in.(type) {
	case order.MapSlice:
		return treeFromMapSlice(in)
	case []interface{}:
		out := make([]interface{}, 0, len(in))
		for _, v := range in {
			out = append(out, jsonToTOML(v))
		}
		return sliceToTrees(out)
	default:
		return tomlFromSimple(in)
	}
}

func sliceToTrees(vs []interface{}) interface{} {
	var out []*gotoml.Tree
	for _, v := range vs {
		t, ok := v.(*gotoml.Tree)
		if !ok {
			return vs
		}
		out = append(out, t)
	}
	if len(out) > 0 {
		return out
	}
	return vs
}

func treeFromMapSlice(m order.MapSlice) *gotoml.Tree {
	tree, err := gotoml.TreeFromMap(map[string]interface{}{})
	if err != nil {
		return nil
	}
	for _, item := range m {
		key, ok := item.Key.(string)
		if !ok {
			continue
		}
		tree.SetPath([]string{key}, jsonToTOML(item.Val))
	}
	return tree
}

func tomlFromSimple(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
	}()
	tree, err := gotoml.TreeFromMap(map[string]interface{}{"v": v})
	if err != nil {
		return v
	}
	return tree.GetPath([]string{"v"})
}

func (TOML) Decode(r io.Reader) (interface{}, error) {
	tree, err := gotoml.LoadReader(r)
	if err != nil {
		return nil, err
	}
	return tomlToJSON(tree), nil
}

func tomlToJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case *gotoml.Tree:
		keys := v.Keys()
		out := make(TOMLTrees, 0, len(keys))
		for _, key := range keys {
			out = append(out, TOMLTree{
				Key: key,
				Val: tomlToJSON(v.GetPath([]string{key})),
				Pos: v.GetPositionPath([]string{key}),
			})
		}
		sort.Sort(out)
		return out.MapSlice()
	case []*gotoml.Tree:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, tomlToJSON(item))
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, tomlToJSON(item))
		}
		return out
	default:
		return tomlToSimple(v)
	}
}

// ensures explicit unmarshaling from tree -- may be unnecessary
func tomlToSimple(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
	}()
	tree, err := gotoml.TreeFromMap(map[string]interface{}{"v": v})
	if err != nil {
		return v
	}
	sMap := map[string][]interface{}{}
	if err := tree.Unmarshal(&sMap); err == nil {
		return sMap["v"]
	}
	vMap := map[string]interface{}{}
	if err := tree.Unmarshal(&vMap); err == nil {
		return vMap["v"]
	}
	return v
}

type TOMLTree struct {
	Key string
	Val interface{}
	Pos gotoml.Position
}

type TOMLTrees []TOMLTree

func (t TOMLTrees) Len() int      { return len(t) }
func (t TOMLTrees) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t TOMLTrees) Less(i, j int) bool {
	if t[i].Pos.Line == t[j].Pos.Line {
		return t[i].Pos.Col < t[j].Pos.Col
	}
	return t[i].Pos.Line < t[j].Pos.Line
}

func (t TOMLTrees) MapSlice() order.MapSlice {
	var out order.MapSlice
	for _, item := range t {
		out = append(out, order.MapItem{
			Key: item.Key,
			Val: item.Val,
		})
	}
	return out
}
