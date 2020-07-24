package convert

import (
	"fmt"
	"io"
	"math"
	"sort"

	gotoml "github.com/pelletier/go-toml"

	"github.com/sclevine/yj/order"
)

type TOML struct {
	FloatStrings bool
	Indent       bool
}

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) (err error) {
	defer catchFailure(&err)
	enc := gotoml.NewEncoder(&trimWriter{w: w})
	enc.Order(gotoml.OrderPreserve)
	if !t.Indent {
		enc.Indentation("")
	}
	converter := tomlConverter{floatStrings: t.FloatStrings}
	return enc.Encode(converter.toTOML(in))
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

type tomlConverter struct {
	line         int
	floatStrings bool
}

func (l *tomlConverter) toTOML(val interface{}) interface{} {
	switch val := val.(type) {
	case order.MapSlice:
		return l.mapToTree(val)
	case []interface{}:
		out := make([]interface{}, 0, len(val))
		for _, v := range val {
			if v == nil {
				continue
			}
			out = append(out, l.toTOML(v))
		}
		return sliceToTrees(out)
	default:
		return l.simpleToTOML(val)
	}
}

func newTree() *gotoml.Tree {
	tree, err := gotoml.TreeFromMap(map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	return tree
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

func tomlKey(s string) string {
	if s == "" {
		return "\"\""
	}
	if len(s) < 2 {
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return `"\` + s[:len(s)-1] + `\""`
	}
	return s
}

func (l *tomlConverter) mapToTree(m order.MapSlice) *gotoml.Tree {
	tree := newTree()
	tl := TreesLast(m)
	sort.Stable(tl)
	for _, item := range tl {
		key, ok := item.Key.(string)
		if !ok {
			panic(fmt.Errorf("non-string key: %#v", item.Key))
		}
		if item.Val == nil {
			continue
		}
		keys := []string{tomlKey(key)}
		line := l.line
		l.line++
		tree.SetPath(keys, l.toTOML(item.Val))
		tree.SetPositionPath(keys, gotoml.Position{Line: line})
	}
	return tree
}

func (l *tomlConverter) simpleToTOML(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
		out = l.denormalize(out)
	}()
	tree, err := gotoml.TreeFromMap(map[string]interface{}{"v": v})
	if err != nil {
		return v
	}
	tree.SetPositionPath([]string{"v"}, gotoml.Position{Line: l.line})
	l.line++
	return tree.GetPath([]string{"v"})
}

func (l *tomlConverter) denormalize(in interface{}) interface{} {
	if l.floatStrings {
		switch in {
		case "NaN":
			return math.NaN()
		case "Infinity":
			return math.Inf(1)
		case "-Infinity":
			return math.Inf(-1)
		}
	}
	return in
}

func (t TOML) Decode(r io.Reader) (out interface{}, err error) {
	defer catchFailure(&err)
	tree, err := gotoml.LoadReader(r)
	if err != nil {
		return nil, err
	}
	return t.tomlToJSON(tree), nil
}

func (t TOML) tomlToJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case *gotoml.Tree:
		keys := v.Keys()
		out := make(TOMLTrees, 0, len(keys))
		for _, key := range keys {
			out = append(out, TOMLTree{
				Key: key,
				Val: t.tomlToJSON(v.GetPath([]string{key})),
				Pos: v.GetPositionPath([]string{key}),
			})
		}
		sort.Sort(out)
		return out.MapSlice()
	case []*gotoml.Tree:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, t.tomlToJSON(item))
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, t.tomlToJSON(item))
		}
		return out
	default:
		return t.tomlToSimple(v)
	}
}

// ensures explicit unmarshaling from tree -- may be unnecessary
func (t TOML) tomlToSimple(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
		out = t.normalize(out)
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

func (t TOML) normalize(in interface{}) interface{} {
	switch in := in.(type) {
	case float64:
		if t.FloatStrings {
			switch {
			case math.IsNaN(in):
				return "NaN"
			case math.IsInf(in, 1):
				return "Infinity"
			case math.IsInf(in, -1):
				return "-Infinity"
			}
		} else {
			switch {
			case math.IsNaN(in):
				return (*float64)(nil)
			case math.IsInf(in, 1):
				return math.MaxFloat64
			case math.IsInf(in, -1):
				return -math.MaxFloat64
			}
		}
	}
	return in
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

type TreesLast order.MapSlice

func (t TreesLast) Len() int      { return len(t) }
func (t TreesLast) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t TreesLast) Less(i, j int) bool {
	return !isMapSlices(t[i].Val) && isMapSlices(t[j].Val)
}

func isMapSlices(v interface{}) bool {
	switch v := v.(type) {
	case order.MapSlice:
		return true
	case []interface{}:
		for _, u := range v {
			if _, ok := u.(order.MapSlice); !ok {
				return false
			}
		}
		return len(v) > 0
	default:
		return false
	}
}
