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

const HelpMsg = `Usage: %s [-][rcjenkh]

Converts stdin from JSON/YAML to YAML/JSON.

-r     Convert JSON to YAML instead of YAML to JSON
-c     Use CandiedYAML parser instead of GoYAML parser
-n     Do not covert infinity, -infinity, and NaN to/from strings
-h     Show this help message

YAML to JSON options:

-e     Escape HTML in JSON output (ignored for JSON to YAML)

JSON to YAML (-r) options:

-j     Use a JSON parser instead of a YAML parser to decode JSON
-k     Attempt to parse keys as JSON objects/numbers
`

func main() {
	config, err := args.Parse(os.Args[1:])
	if err != nil {
		fail(err)
	}
	if config.Help {
		fmt.Printf(HelpMsg, os.Args[0])
		os.Exit(0)
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

	decoder.Unmarshal = goyaml.Unmarshal
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
	encoder := &yaml.Encoder{}
	if config.JSONKeys {
		encoder.KeyUnmarshal = (&yaml.JSON{JSONDecoder: config.JSONDecoder}).Unmarshal
	}

	if config.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}

	encoder.Marshal = goyaml.Marshal
	if config.CandiedYAML {
		encoder.Marshal = candiedyaml.Marshal
	}

	unmarshalFunc := goyaml.Unmarshal
	if config.JSONDecoder {
		unmarshalFunc = json.Unmarshal
	}

	var data interface{}
	if err := unmarshalFunc(input, &data); err != nil {
		return nil, err
	}

	return encoder.YAML(data)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	fmt.Printf(HelpMsg, os.Args[0])
	os.Exit(1)
}
