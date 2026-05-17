package main

import (
	"fmt"
	"io"
	"os"
)

const HelpMsg = `Usage: %s [-][ytjcrneikhv]

Convert between YAML, TOML, JSON, and HCL.
Preserves map order.

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
-n     Do not convert inf, -inf, and NaN to/from strings (YAML or TOML only)
-e     Escape HTML (JSON out only)
-i     Indent output (JSON or TOML out only)
-k     Attempt to parse keys as objects or numeric types (YAML out only)
-h     Show this help message
-v     Show version

`

var Version = "0.0.0"

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
	if config.Version {
		fmt.Fprintln(stdout, "v"+Version)
		return 0
	}

	rep, err := config.From.Decode(stdin)
	if err != nil {
		// TODO: I forget if there's a reason this isn't io.EOF
		if err.Error() == "EOF" {
			return 0
		}
		fmt.Fprintf(stderr, "Error parsing %s: %s\n", config.From, err)
		return 1
	}
	if err := config.To.Encode(stdout, rep); err != nil {
		fmt.Fprintf(stderr, "Error writing %s: %s\n", config.To, err)
		return 1
	}
	return 0
}
