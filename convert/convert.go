package convert

import (
	"bytes"
	"encoding/json"
	"math"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	hcljson "github.com/hashicorp/hcl/json/parser"
	goyaml "gopkg.in/yaml.v2"

	"github.com/sclevine/yj/yaml"
)

type (
	YAML struct {
		FloatStrings bool
		JSONKeys     bool
		EscapeHTML   bool
	}
	TOML struct{}
	JSON struct {
		EscapeHTML bool
	}
	HCL struct{}
)

type Encoding interface {
	String() string
	Encode(input interface{}) ([]byte, error)
	Decode(input []byte) (interface{}, error)
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

func (TOML) String() string {
	return "TOML"
}

func (TOML) Encode(input interface{}) ([]byte, error) {
	output := &bytes.Buffer{}
	err := toml.NewEncoder(output).Encode(input)
	return output.Bytes(), err
}

func (TOML) Decode(input []byte) (interface{}, error) {
	var data interface{}
	return data, toml.Unmarshal(input, &data)
}

func (JSON) String() string {
	return "JSON"
}

func (j JSON) Encode(input interface{}) ([]byte, error) {
	output := &bytes.Buffer{}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(j.EscapeHTML)
	err := encoder.Encode(input)
	return output.Bytes(), err
}

func (JSON) Decode(input []byte) (interface{}, error) {
	var data interface{}
	return data, json.Unmarshal(input, &data)
}

func (HCL) String() string {
	return "HCL"
}

func (HCL) Encode(input interface{}) ([]byte, error) {
	json, err := JSON{}.Encode(input)
	if err != nil {
		return nil, err
	}

	ast, err := hcljson.Parse(json)
	if err != nil {
		return nil, err
	}

	output := &bytes.Buffer{}
	err = printer.Fprint(output, ast)
	return output.Bytes(), err
}

func (HCL) Decode(input []byte) (interface{}, error) {
	var data interface{}
	return data, hcl.Unmarshal(input, &data)
}
