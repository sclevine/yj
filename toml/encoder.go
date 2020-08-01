package toml

import (
	"fmt"
	"math"
	"sort"

	gotoml "github.com/pelletier/go-toml"

	"github.com/sclevine/yj/order"
)

type Encoder struct {
	// If set, the set values will be converted to NaN, Inf, etc.
	NaN, PosInf, NegInf interface{}
}

func (e *Encoder) Encode(normal interface{}) (toml interface{}, err error) {
	defer catchFailure(&err)
	et := encodeTracker{Encoder: e}
	return et.denormalize(normal), nil
}

type encodeTracker struct {
	*Encoder
	line int
}

func (e *encodeTracker) denormalize(val interface{}) interface{} {
	switch val := val.(type) {
	case order.MapSlice:
		return e.mapToTree(val)
	case []interface{}:
		out := make([]interface{}, 0, len(val))
		for _, v := range val {
			if v == nil {
				continue
			}
			out = append(out, e.denormalize(v))
		}
		return sliceToTrees(out)
	default:
		return e.simpleToTOML(val)
	}
}

func (e *encodeTracker) mapToTree(m order.MapSlice) *gotoml.Tree {
	tree := newTree()
	tl := treesLast(m)
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
		line := e.line
		e.line++
		tree.SetPath(keys, e.denormalize(item.Val))
		tree.SetPositionPath(keys, gotoml.Position{Line: line})
	}
	return tree
}

func (e *encodeTracker) simpleToTOML(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
		out = e.postprocess(out)
	}()
	tree, err := gotoml.TreeFromMap(map[string]interface{}{"v": v})
	if err != nil {
		return v
	}
	tree.SetPositionPath([]string{"v"}, gotoml.Position{Line: e.line})
	e.line++
	return tree.GetPath([]string{"v"})
}

func (e *encodeTracker) postprocess(in interface{}) interface{} {
	switch in {
	case nil:
		return nil
	case e.NaN:
		return math.NaN()
	case e.PosInf:
		return math.Inf(1)
	case e.NegInf:
		return math.Inf(-1)
	}
	return in
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
