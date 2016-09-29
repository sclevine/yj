package yaml

import (
	"fmt"
	"math"
	"reflect"
)

type Decoder struct {
	Unmarshal  func([]byte, interface{}) error
	KeyMarshal func(interface{}) ([]byte, error)

	// If not set, input YAML must not contain these.
	// These are returned unmodified in the output of JSON.
	NaN, PosInf, NegInf interface{}
}

// JSON decodes YAML into an object that marshals cleanly into JSON.
func (d *Decoder) JSON(yaml []byte) (json interface{}, err error) {
	defer catchFailure(&err)

	// Must pass *interface{} due to go-yaml quirk
	var data interface{}
	if err := d.Unmarshal(yaml, &data); err != nil {
		return nil, err
	}
	return d.jsonify(data), nil
}

func (d *Decoder) jsonify(in interface{}) interface{} {
	switch in := in.(type) {
	case map[interface{}]interface{}:
		out := map[string]interface{}{}
		for k, v := range in {
			out[d.jsonifyKey(k)] = d.jsonify(v)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(in))
		for i, v := range in {
			out[i] = d.jsonify(v)
		}
		return out
	case float64:
		return d.jsonifyFloat(in)
	default:
		return d.jsonifyOther(in)
	}
}

func (d *Decoder) jsonifyOther(in interface{}) interface{} {
	switch reflect.ValueOf(in).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Float32:
		panic(fmt.Errorf("unexpected type: %#v", in))
	}
	return in
}

func (d *Decoder) jsonifyFloat(in float64) interface{} {
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

func (d *Decoder) jsonifyKey(in interface{}) string {
	switch key := d.jsonify(in).(type) {
	case string:
		return key
	case fmt.Stringer:
		return key.String()
	default:
		out, err := d.KeyMarshal(key)
		if err != nil {
			panic(err)
		}
		return string(out)

	}
}
