package yaml

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/sclevine/yj/order"
)

type Encoder struct {
	KeyUnmarshal func([]byte, interface{}) error

	// If set, the set values will be converted to NaN, Inf, etc.
	NaN, PosInf, NegInf          interface{}
	KeyNaN, KeyPosInf, KeyNegInf interface{}
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
	case reflect.Float64:
		out := strconv.FormatFloat(in.(float64), 'g', -1, 64)
		if isNumbers(out) {
			return &yaml.Node{Kind: yaml.ScalarNode, Value: out + ".0"}
		}
	}
	return in
}

func isNumbers(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func (e *Encoder) key(in string) interface{} {
	var key interface{} = in
	if e.KeyUnmarshal != nil {
		var v interface{}
		if err := e.KeyUnmarshal([]byte(in), &v); err == nil {
			key = v
		}
	}
	kenc := *e
	kenc.NaN = e.KeyNaN
	kenc.PosInf = e.KeyPosInf
	kenc.NegInf = e.KeyNegInf
	switch out := kenc.denormalize(key); reflect.ValueOf(out).Kind() {
	case reflect.Map, reflect.Slice, reflect.Func:
		return &out
	default:
		return out
	}
}
