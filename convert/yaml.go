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
	encoder := &yaml.Encoder{
		EncodeYAML: func(w io.Writer, v interface{}) error {
			enc := goyaml.NewEncoder(w)
			enc.SetIndent(2)
			return enc.Encode(v)
		},
	}
	if y.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}
	if y.JSONKeys {
		encoder.KeyUnmarshal = (&yaml.JSON{}).Unmarshal
	}
	return encoder.YAML(w, in)
}

func (y YAML) Decode(r io.Reader) (interface{}, error) {
	decoder := &yaml.Decoder{
		DecodeYAML: func(r io.Reader) (*goyaml.Node, error) {
			var data goyaml.Node
			return &data, goyaml.NewDecoder(r).Decode(&data)
		},
		KeyMarshal: (&yaml.JSON{EscapeHTML: y.EscapeHTML}).Marshal, // FIXME: double-check map-keys

		NaN:    (*float64)(nil),
		PosInf: math.MaxFloat64,
		NegInf: -math.MaxFloat64,
	}

	if y.FloatStrings {
		decoder.NaN = "NaN"
		decoder.PosInf = "Infinity"
		decoder.NegInf = "-Infinity"
	}
	return decoder.JSON(r)
}
