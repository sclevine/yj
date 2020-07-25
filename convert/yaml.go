package convert

import (
	"io"

	goyaml "gopkg.in/yaml.v3"

	"github.com/sclevine/yj/yaml"
)

type YAML struct {
	SpecialFloats
	KeySpecialFloats SpecialFloats
	JSONKeys         bool
	EscapeHTML       bool
}

func (YAML) String() string {
	return "YAML"
}

func (y YAML) Encode(w io.Writer, in interface{}) error {
	enc := &yaml.Encoder{
		NaN:       y.NaN(),
		PosInf:    y.PosInf(),
		NegInf:    y.NegInf(),
		KeyNaN:    y.KeySpecialFloats.NaN(),
		KeyPosInf: y.KeySpecialFloats.PosInf(),
		KeyNegInf: y.KeySpecialFloats.NegInf(),
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
		NaN:        y.NaN(),
		PosInf:     y.PosInf(),
		NegInf:     y.NegInf(),
		KeyNaN:     y.KeySpecialFloats.NaN(),
		KeyPosInf:  y.KeySpecialFloats.PosInf(),
		KeyNegInf:  y.KeySpecialFloats.NegInf(),
	}
	return dec.Decode(&node)
}
