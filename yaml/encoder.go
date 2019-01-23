package yaml

import (
	"fmt"
	"math"
	"reflect"
)

type Encoder struct {
	EncodeYAML   func(interface{}) error
	KeyUnmarshal func([]byte, interface{}) error

	// If set, these will be converted to floats.
	NaN, PosInf, NegInf interface{}
}

// YAML encodes a JSON object into YAML.
// Special string keys from the Decoder are accounted for.
// YAML objects are accepted, as long as represent valid JSON.
// Internal structs are currently passed through unmodified.
func (e *Encoder) YAML(json interface{}) (err error) {
	defer catchFailure(&err)
	return e.EncodeYAML(e.yamlify(json))
}

func (e *Encoder) yamlify(in interface{}) interface{} {
	switch in := in.(type) {
	case map[string]interface{}:
		out := map[interface{}]interface{}{}
		for k, v := range in {
			out[e.yamlifyKey(k)] = e.yamlify(v)
		}
		return out
	case map[interface{}]interface{}: // TODO: test
		out := map[interface{}]interface{}{}
		for k, v := range in {
			switch k := k.(type) {
			case string:
				out[e.yamlifyKey(k)] = e.yamlify(v)
			default:
				panic(fmt.Errorf("invalid key: %#v", k)) // test!
			}
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(in))
		for i, v := range in {
			out[i] = e.yamlify(v)
		}
		return out
	case []map[string]interface{}: // TODO: test
		out := make([]interface{}, len(in))
		for i, v := range in {
			out[i] = e.yamlify(v)
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
