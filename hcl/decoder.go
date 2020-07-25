package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"github.com/sclevine/yj/order"
)

type Decoder struct{}

func (d *Decoder) Decode(node ast.Node) (out interface{}, err error) {
	defer catchFailure(&err)
	return d.normalize(node), nil
}

func (d *Decoder) normalize(node ast.Node) interface{} {
	switch n := node.(type) {
	case *ast.ObjectList:
		done := make(map[string]struct{})
		list := newObjList()
		for _, item := range n.Items {
			if item.Val == nil {
				continue
			}
			if len(item.Keys) == 0 {
				panic(fmt.Errorf("empty key at line %d", item.Pos().Line))
			}
			key := item.Keys[0].Token.Value().(string)
			if _, ok := done[key]; ok {
				continue
			}
			itemVal := item.Val
			if len(item.Keys) > 1 {
				itemVal = n.Filter(key)
				done[key] = struct{}{}
			}
			val := d.normalize(itemVal)
			switch val := val.(type) {
			case []interface{}:
				for _, v := range val {
					list.append(key, v)
				}
			case order.MapSlice:
				list.append(key, val)
			default:
				list.set(key, val)
			}
		}
		return list.mapSlice()
	case *ast.ObjectType:
		return d.normalize(n.List)
	case *ast.ListType:
		out := make([]interface{}, 0, len(n.List))
		for _, item := range n.List {
			out = append(out, d.normalize(item))
		}
		return out
	case *ast.LiteralType:
		var out interface{}
		if err := hcl.DecodeObject(&out, n); err != nil {
			panic(err)
		}
		return out
	default:
		panic(fmt.Errorf("invalid type: %#v", n))
	}
}

type objList struct {
	keys []string
	m    map[string]interface{}
}

func newObjList() *objList {
	return &objList{
		m: make(map[string]interface{}),
	}
}

func (l *objList) append(k string, v interface{}) {
	if dst, ok := l.m[k]; ok {
		dstList, ok := dst.([]interface{})
		if !ok {
			panic(fmt.Errorf("cannot append '%#v' to '%#v", v, dst))
		}
		l.m[k] = append(dstList, v)
	} else {
		l.keys = append(l.keys, k)
		l.m[k] = []interface{}{v}
	}
}

func (l *objList) set(k string, v interface{}) {
	if dst, ok := l.m[k]; ok {
		if dstList, ok := dst.([]interface{}); ok {
			panic(fmt.Errorf("cannot set '%#v' over '%#v", v, dstList))
		}
	} else {
		l.keys = append(l.keys, k)
	}
	l.m[k] = v
}

func (l *objList) mapSlice() order.MapSlice {
	out := make(order.MapSlice, 0, len(l.keys))
	for _, key := range l.keys {
		out = append(out, order.MapItem{
			Key: key,
			Val: l.m[key],
		})
	}
	return out
}

func catchFailure(err *error) {
	if r := recover(); r != nil {
		var ok bool
		if *err, ok = r.(error); !ok {
			*err = fmt.Errorf("unexpected failure: %v", r)
		}
	}
}
