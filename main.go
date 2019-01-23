package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
          -cj, -c = HCL to JSON
          -cc     = HCL to HCL
-n     Do not covert inf, -inf, and NaN to/from strings (YAML in/out only)
-e     Escape HTML (JSON out only)
-k     Attempt to parse keys as objects or numbers types (YAML out only)
-h     Show this help message

`

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, osArgs []string) (code int) {
	config, err := Parse(osArgs[1:]...)
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

	if len(bytes.TrimSpace(input)) == 0 {
		return 0
	}

	// TODO: if from == to, don't do yaml decode/encode to avoid stringifying the keys
	rep, err := config.From.Decode(input)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing %s: %s\n", config.From, err)
		return 1
	}
	output, err := config.To.Encode(rep)
	if err != nil {
		fmt.Fprintf(stderr, "Error writing %s: %s\n", config.To, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s", output)
	return 0
}
