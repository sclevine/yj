package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
-n     Do not covert Infinity, -Infinity, and NaN to/from strings
-h     Show this help message

YAML to JSON options:

-e     Escape HTML in JSON output (ignored for JSON to YAML)

JSON to YAML (-r) options:

-y     Use a YAML decoder instead of a JSON decoder to parse JSON
-k     Attempt to parse keys as JSON objects/numbers
`

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, osArgs []string) (code int) {
	config, err := args.Parse(osArgs[1:]...)
	if err != nil {
		failMsg(stderr, err)
		return 1
	}
	if config.Help {
		fmt.Fprintf(stdout, HelpMsg, os.Args[0])
		return 0
	}

	input, err := ioutil.ReadAll(stdin)
	if err != nil {
		failMsg(stderr, err)
		return 1
	}

	convertFunc := convertYAMLToJSON
	if config.Reverse {
		convertFunc = convertJSONToYAML
	}

	output, err := convertFunc(input, config)
	if err != nil {
		failMsg(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s", output)

	return 0
}

func convertYAMLToJSON(input []byte, config *args.Config) ([]byte, error) {
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

func convertJSONToYAML(input []byte, config *args.Config) ([]byte, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	encoder := &yaml.Encoder{}

	if config.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}

	encoder.Marshal = goyaml.Marshal
	if config.CandiedYAML {
		encoder.Marshal = candiedyaml.Marshal
	}

	unmarshalFunc := json.Unmarshal
	if config.JSONAsYAML {
		unmarshalFunc = goyaml.Unmarshal
		if config.CandiedYAML {
			unmarshalFunc = candiedyaml.Unmarshal
		}
	}

	if config.JSONKeys {
		encoder.KeyUnmarshal = unmarshalFunc
	}

	var data interface{}
	if err := unmarshalFunc(input, &data); err != nil {
		return nil, err
	}

	return encoder.YAML(data)
}

func failMsg(out io.Writer, err error) {
	fmt.Fprintf(out, "Error: %s\n", err)
	fmt.Fprintf(out, HelpMsg, os.Args[0])
}
