package convert

import (
	"io"
	"math"

	goyaml "gopkg.in/yaml.v3"

	"github.com/sclevine/yj/yaml"
)

type YAML struct {
	FloatStrings bool
	JSONKeys     bool
	EscapeHTML   bool
}

func (YAML) String() string {
	return "YAML"
}

func (y YAML) Encode(w io.Writer, in interface{}) error {
	enc := &yaml.Encoder{}
	if y.FloatStrings {
		enc.NaN = "NaN"
		enc.PosInf = "Infinity"
		enc.NegInf = "-Infinity"
	}
	if y.JSONKeys {
		enc.KeyUnmarshal = (&yaml.KeyJSON{}).Unmarshal
	}
	out, err := enc.Encode(in)
	if err != nil {
		return err
	}
	yamlEnc := goyaml.NewEncoder(w)
	yamlEnc.SetIndent(2)
	return yamlEnc.Encode(out)
}

func (y YAML) Decode(r io.Reader) (interface{}, error) {
	var node goyaml.Node
	if err := goyaml.NewDecoder(r).Decode(&node); err != nil {
		return nil, err
	}
	dec := &yaml.Decoder{
		KeyMarshal: (&yaml.KeyJSON{EscapeHTML: y.EscapeHTML}).Marshal, // FIXME: double-check map-keys
		NaN:    (*float64)(nil),
		PosInf: math.MaxFloat64,
		NegInf: -math.MaxFloat64,
	}
	if y.FloatStrings {
		dec.NaN = "NaN"
		dec.PosInf = "Infinity"
		dec.NegInf = "-Infinity"
	}
	return dec.Decode(&node)
}
