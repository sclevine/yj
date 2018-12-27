package args

import (
	"fmt"
	"strings"
)

type Config struct {
	From, To     Encoding
	EscapeHTML   bool
	FloatStrings bool
	JSONKeys     bool
	Help         bool
}

type Encoding rune

const (
	YAML Encoding = FlagYAML
	TOML Encoding = FlagTOML
	JSON Encoding = FlagJSON
	HCL  Encoding = FlagHCL

	FlagYAML           = 'y'
	FlagTOML           = 't'
	FlagJSON           = 'j'
	FlagHCL            = 'c'
	FlagReverse        = 'r'
	FlagNoFloatStrings = 'n'
	FlagEscapeHTML     = 'e'
	FlagJSONKeys       = 'k'
	FlagHelp           = 'h'
)

func Parse(args ...string) (*Config, error) {
	flatArgs := strings.Join(args, "")
	invalidArgs := strings.Split(strings.Map(flagFilter, flatArgs), "")

	if len(invalidArgs) > 0 {
		return nil, fmt.Errorf("invalid flags specified: %s", strings.Join(invalidArgs, " "))
	}

	config := &Config{
		EscapeHTML:   strings.ContainsRune(flatArgs, FlagEscapeHTML),
		FloatStrings: !strings.ContainsRune(flatArgs, FlagNoFloatStrings),
		JSONKeys:     strings.ContainsRune(flatArgs, FlagJSONKeys),
		Help:         strings.ContainsRune(flatArgs, FlagHelp),
	}
	config.From, config.To = transform(flatArgs)

	if config.JSONKeys && config.To != YAML {
		return nil, fmt.Errorf("flag -%c only valid for YAML output", FlagJSONKeys)
	}
	if config.EscapeHTML && config.To != JSON {
		return nil, fmt.Errorf("flag -%c only valid for JSON output", FlagEscapeHTML)
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagYAML, FlagTOML, FlagJSON, FlagHCL, FlagReverse,
		FlagEscapeHTML, FlagNoFloatStrings, FlagJSONKeys,
		FlagHelp, '\t', ' ', '-':
		return -1
	}
	return r
}

func transform(s string) (from, to Encoding) {
	for _, r := range s {
		switch r {
		case FlagYAML, FlagTOML, FlagJSON, FlagHCL:
			from, to = to, Encoding(r)
		case FlagReverse:
			from, to = JSON, YAML
		}
	}
	if from == 0 {
		if to == 0 {
			to = YAML
		}
		from, to = to, JSON
	}
	return
}
