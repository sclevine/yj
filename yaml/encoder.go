package yaml

import (
	"fmt"
	"math"
	"reflect"

	"github.com/sclevine/yj/order"
)

type Encoder struct {
	KeyUnmarshal func([]byte, interface{}) error

	// If set, these will be converted to floats.
	NaN, PosInf, NegInf interface{}
}

// Encode encodes the normalized object format into a suitable format for marshaling to YAML.
func (e *Encoder) Encode(normal interface{}) (yaml interface{}, err error) {
	defer catchFailure(&err)
	return e.denormalize(normal), nil
}

func (e *Encoder) denormalize(in interface{}) interface{} {
	switch in := in.(type) {
	case order.MapSlice:
		out := make(order.MapSlice, 0, len(in))
		for _, item := range in {
			key, ok := item.Key.(string)
			if !ok {
				panic(fmt.Errorf("key not string: %#v", item.Key))
			}
			out = append(out, order.MapItem{
				Key: e.key(key),
				Val: e.denormalize(item.Val),
			})
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(in))
		for _, v := range in {
			out = append(out, e.denormalize(v))
		}
		return out
	default:
		return e.other(in)
	}
}

func (e *Encoder) other(in interface{}) interface{} {
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
	switch reflect.ValueOf(in).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Float32:
		panic(fmt.Errorf("unexpected type: %#v", in))
	}
	return in
}

func (e *Encoder) key(in string) interface{} {
	var key interface{} = in
	if e.KeyUnmarshal != nil {
		var v interface{}
		if err := e.KeyUnmarshal([]byte(in), &v); err == nil {
			key = v
		}
	}
	switch out := e.denormalize(key); reflect.ValueOf(out).Kind() {
	case reflect.Map, reflect.Slice, reflect.Func:
		return &out
	default:
		return out
	}
}
