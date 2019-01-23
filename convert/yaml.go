package convert

import (
	"encoding/json"
	"math"

	goyaml "gopkg.in/yaml.v2"

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

func (y YAML) Encode(input interface{}) ([]byte, error) {
	encoder := &yaml.Encoder{}
	if y.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}
	encoder.Marshal = goyaml.Marshal
	if y.JSONKeys {
		encoder.KeyUnmarshal = json.Unmarshal
	}
	return encoder.YAML(input)
}

func (y YAML) Decode(input []byte) (interface{}, error) {
	decoder := &yaml.Decoder{
		KeyMarshal: (&yaml.JSON{EscapeHTML: y.EscapeHTML}).Marshal,

		NaN:    (*float64)(nil),
		PosInf: math.MaxFloat64,
		NegInf: -math.MaxFloat64,
	}

	if y.FloatStrings {
		decoder.NaN = "NaN"
		decoder.PosInf = "Infinity"
		decoder.NegInf = "-Infinity"
	}

	decoder.Unmarshal = goyaml.Unmarshal
	return decoder.JSON(input)
}
