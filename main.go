package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	goyaml "gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
	"github.com/sclevine/yj/args"
	"github.com/sclevine/yj/yaml"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
)

const HelpMsg = `Usage: %s [-][ytjcrnekh]

Convert YAML, TOML, JSON, or HCL to YAML, TOML, JSON, or HCL.

-x[x]  Convert using stdin. Valid options:
          -yj, -y = YAML to JSON (default)
          -yy     = YAML to YAML
          -yt     = YAML to TOML
          -yc     = YAML to HCL
          -tj, -t = TOML to JSON
          -ty     = TOML to YAML
          -tt     = TOML to TOML
          -tc     = TOML to HCL
          -jj     = JSON to JSON
          -jy, -r = JSON to YAML
          -jt     = JSON to TOML
          -jc     = JSON to HCL
          -cy     = HCL to YAML
          -ct     = HCL to TOML
          -cj  -c = HCL to JSON
          -cc     = HCL to HCL
-n     Do not covert Infinity, -Infinity, and NaN to/from strings
-e     Escape HTML (JSON output only)
-k     Attempt to parse keys as objects or numbers types (YAML output only)
-h     Show this help message

`

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, osArgs []string) (code int) {
	config, err := args.Parse(osArgs[1:]...)
	if err != nil {
		fmt.Fprintf(stderr, HelpMsg, os.Args[0])
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}
	if config.Help {
		fmt.Fprintf(stdout, HelpMsg, os.Args[0])
		return 0
	}

	input, err := ioutil.ReadAll(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}

	var from func([]byte, *args.Config) (interface{}, error)
	switch config.From {
	case args.YAML:
		from = fromYAML
	case args.TOML:
		from = fromTOML
	case args.JSON:
		from = fromJSON
	case args.HCL:
		from = fromHCL
	}

	var to func(interface{}, *args.Config) ([]byte, error)
	switch config.To {
	case args.YAML:
		to = toYAML
	case args.TOML:
		to = toTOML
	case args.JSON:
		to = toJSON
	case args.HCL:
		to = toHCL
	}

	// TODO: if from == to, don't do yaml decode/encode to avoid stringifying the keys
	rep, err := from(input, config)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing %s: %s\n", config.From, err)
		return 1
	}
	output, err := to(rep, config)
	if err != nil {
		fmt.Fprintf(stderr, "Error writing %s: %s\n", config.To, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s", output)
	return 0
}

func fromYAML(input []byte, config *args.Config) (interface{}, error) {
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
	return decoder.JSON(input)
}

func fromTOML(input []byte, _ *args.Config) (interface{}, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	var data interface{}
	return data, toml.Unmarshal(input, &data)
}

func fromJSON(input []byte, _ *args.Config) (interface{}, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	var data interface{}
	return data, json.Unmarshal(input, &data)
}

func fromHCL(input []byte, _ *args.Config) (interface{}, error) {
	if len(bytes.TrimSpace(input)) == 0 {
		return nil, nil
	}
	var data interface{}
	return data, hcl.Unmarshal(input, &data)
}

func toYAML(input interface{}, config *args.Config) ([]byte, error) {
	encoder := &yaml.Encoder{}
	if config.FloatStrings {
		encoder.NaN = "NaN"
		encoder.PosInf = "Infinity"
		encoder.NegInf = "-Infinity"
	}
	encoder.Marshal = goyaml.Marshal
	if config.JSONKeys {
		encoder.KeyUnmarshal = json.Unmarshal
	}
	return encoder.YAML(input)
}

func toTOML(input interface{}, _ *args.Config) ([]byte, error) {
	output := &bytes.Buffer{}
	err := toml.NewEncoder(output).Encode(input)
	return output.Bytes(), err
}

func toJSON(input interface{}, config *args.Config) ([]byte, error) {
	output := &bytes.Buffer{}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(config.EscapeHTML)
	err := encoder.Encode(input)
	return output.Bytes(), err
}

func toHCL(input interface{}, config *args.Config) ([]byte, error) {
	json, err := toJSON(input, config)
	if err != nil {
		return nil, err
	}

	ast, err := jsonParser.Parse(json)
	if err != nil {
		return nil, err
	}

	output := &bytes.Buffer{}
	err = printer.Fprint(output, ast)
	return output.Bytes(), err
}
