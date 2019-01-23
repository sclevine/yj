package convert

import (
	"bytes"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	hcljson "github.com/hashicorp/hcl/json/parser"
)

type HCL struct{}

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
