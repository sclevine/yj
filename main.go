package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
	goyaml "gopkg.in/yaml.v2"

	"github.com/sclevine/yj/args"
	"github.com/sclevine/yj/yaml"
)

func main() {
	config, err := args.Parse(os.Args[1:])
	if err != nil {
		fail(err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fail(err)
	}

	convertFunc := ConvertYAMLToJSON
	if config.Reverse {
		convertFunc = ConvertJSONToYAML
	}

	output, err := convertFunc(input, config)
	if err != nil {
		fail(err)
	}
	fmt.Printf("%s", output)
}

func ConvertYAMLToJSON(input []byte, config args.Config) ([]byte, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	decoder := &yaml.Decoder{
		Unmarshal:  goyaml.Unmarshal,
		KeyMarshal: (&yaml.JSON{EscapeHTML: config.EscapeHTML}).Marshal,

		NaN:    (*float64)(nil),
		PosInf: math.MaxFloat64,
		NegInf: -math.MaxFloat64,
	}

	if config.FloatStrings {
		decoder.NaN = "NaN"
		decoder.PosInf = "Infinity"
		decoder.NegInf = "-Infinity"
	}

	if config.CandiedYAML {
		decoder.Unmarshal = candiedyaml.Unmarshal
	}

	data, err := decoder.JSON(input)
	if err != nil {
		return nil, err
	}
	output := &bytes.Buffer{}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(config.EscapeHTML)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func ConvertJSONToYAML(input []byte, config args.Config) ([]byte, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	encoder := &yaml.Encoder{Marshal: goyaml.Marshal}
	if config.JSONKeys {
		encoder.KeyUnmarshal = (&yaml.JSON{EscapeHTML: config.EscapeHTML}).Unmarshal
	}

	if config.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}

	if config.CandiedYAML {
		encoder.Marshal = candiedyaml.Marshal
	}
	var data interface{}
	if err := goyaml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	return encoder.YAML(data)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
