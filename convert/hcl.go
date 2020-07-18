package convert

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
	hcljson "github.com/hashicorp/hcl/json/parser"

	"github.com/sclevine/yj/order"
)

type HCL struct{}

func (HCL) String() string {
	return "HCL"
}

func (HCL) Encode(w io.Writer, in interface{}) error {
	j := &bytes.Buffer{}
	if err := (JSON{}).Encode(j, in); err != nil {
		return err
	}
	f, err := hcljson.Parse(j.Bytes())
	if err != nil {
		return err
	}
	return printer.Fprint(w, f)
}

func (HCL) Decode(r io.Reader) (out interface{}, err error) {
	defer catchFailure(&err)
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f, err := hcl.ParseBytes(in)
	if err != nil {
		return nil, err
	}
	return decode(f.Node), nil
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

// if key is set to literal [LiteralType], only allow literal overrides
// if key is set to list, only allow setting (appending) lists [ListType] and objects [ObjectType] (treated as single-object lists)
// maybe check decoded value instead of checking ast types?

func decode(node ast.Node) interface{} {
	switch n := node.(type) {
	case *ast.ObjectList:
		list := newObjList()
		for _, item := range n.Items {
			key := item.Keys[0].Token.Value().(string)
			val := decode(item.Val)
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
		return decode(n.List)
	case *ast.ListType:
		out := make([]interface{}, 0, len(n.List))
		for _, item := range n.List {
			out = append(out, decode(item))
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

func catchFailure(err *error) {
	if r := recover(); r != nil {
		var ok bool
		if *err, ok = r.(error); !ok {
			*err = fmt.Errorf("unexpected failure: %v", r)
		}
	}
}
