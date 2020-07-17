package yaml

import (
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/sclevine/yj/order"
)

type Encoder struct {
	EncodeYAML   func(io.Writer, interface{}) error
	KeyUnmarshal func([]byte, interface{}) error

	// If set, these will be converted to floats.
	NaN, PosInf, NegInf interface{}
}

// YAML encodes a JSON object into YAML.
// Special string keys from the Decoder are accounted for.
// YAML objects are accepted, as long as represent valid JSON.
// Internal structs are currently passed through unmodified.
func (e *Encoder) YAML(w io.Writer, json interface{}) (err error) {
	defer catchFailure(&err)
	return e.EncodeYAML(w, e.yamlify(json))
}

func (e *Encoder) yamlify(in interface{}) interface{} {
	switch in := in.(type) {
	case order.MapSlice:
		out := make(order.MapSlice, 0, len(in))
		for _, item := range in {
			key, ok := item.Key.(string)
			if !ok {
				panic(fmt.Errorf("key not string: %#v", item.Key))
			}
			out = append(out, order.MapItem{
				Key: e.yamlifyKey(key),
				Val: e.yamlify(item.Val),
			})
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(in))
		for _, v := range in {
			out = append(out, e.yamlify(v))
		}
		return out
	default:
		return e.yamlifyOther(in)
	}
}

func (e *Encoder) yamlifyOther(in interface{}) interface{} {
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

func (e *Encoder) yamlifyKey(in string) interface{} {
	var key interface{} = in
	if e.KeyUnmarshal != nil {
		var v interface{}
		if err := e.KeyUnmarshal([]byte(in), &v); err == nil {
			key = v
		}
	}
	switch out := e.yamlify(key); reflect.ValueOf(out).Kind() {
	case reflect.Map, reflect.Slice, reflect.Func:
		return &out
	default:
		return out
	}
}
