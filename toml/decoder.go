package toml

import (
	"math"
	"sort"

	gotoml "github.com/pelletier/go-toml"

	"github.com/sclevine/yj/order"
)

type Decoder struct {
	// If set, NaN, Inf, etc. are replaced by the set values
	NaN, PosInf, NegInf interface{}
}

func (d *Decoder) Decode(toml interface{}) (normal interface{}, err error) {
	defer catchFailure(&err)
	return d.normalize(toml), nil
}

func (d Decoder) normalize(v interface{}) interface{} {
	switch v := v.(type) {
	case *gotoml.Tree:
		keys := v.Keys()
		out := make(tomlTrees, 0, len(keys))
		for _, key := range keys {
			out = append(out, tomlTree{
				key: key,
				val: d.normalize(v.GetPath([]string{key})),
				pos: v.GetPositionPath([]string{key}),
			})
		}
		sort.Sort(out)
		return out.mapSlice()
	case []*gotoml.Tree:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, d.normalize(item))
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			out = append(out, d.normalize(item))
		}
		return out
	default:
		return d.tomlToSimple(v)
	}
}

// ensures explicit unmarshaling from tree -- may be unnecessary
func (d Decoder) tomlToSimple(v interface{}) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = v
		}
		out = d.postprocess(out)
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

func (d Decoder) postprocess(in interface{}) interface{} {
	switch in := in.(type) {
	case float64:
		switch {
		case d.NaN != nil && math.IsNaN(in):
			return d.NaN
		case d.PosInf != nil && math.IsInf(in, 1):
			return d.PosInf
		case d.NegInf != nil && math.IsInf(in, -1):
			return d.NegInf
		}
		return in
	}
	return in
}

type tomlTree struct {
	key string
	val interface{}
	pos gotoml.Position
}

type tomlTrees []tomlTree

func (t tomlTrees) Len() int      { return len(t) }
func (t tomlTrees) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t tomlTrees) Less(i, j int) bool {
	if t[i].pos.Line == t[j].pos.Line {
		return t[i].pos.Col < t[j].pos.Col
	}
	return t[i].pos.Line < t[j].pos.Line
}

func (t tomlTrees) mapSlice() order.MapSlice {
	var out order.MapSlice
	for _, item := range t {
		out = append(out, order.MapItem{
			Key: item.key,
			Val: item.val,
		})
	}
	return out
}

type treesLast order.MapSlice

func (t treesLast) Len() int      { return len(t) }
func (t treesLast) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t treesLast) Less(i, j int) bool {
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
