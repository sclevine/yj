package args

import (
	"fmt"
	"strings"
)

type Config struct {
	Reverse      bool
	CandiedYAML  bool
	JSONDecoder  bool
	EscapeHTML   bool
	FloatStrings bool
	JSONKeys     bool
	Help         bool
}

const (
	FlagReverse        = 'r'
	FlagCandiedYAML    = 'c'
	FlagJSONDecoder    = 'j'
	FlagEscapeHTML     = 'e'
	FlagNoFloatStrings = 'n'
	FlagJSONKeys       = 'k'
	FlagHelp           = 'h'
)

func Parse(args []string) (Config, error) {
	flatArgs := strings.Join(args, "")

	invalidArgs := strings.Split(strings.Map(flagFilter, flatArgs), "")

	if len(invalidArgs) > 0 {
		return Config{}, fmt.Errorf("invalid flags specified: %s", strings.Join(invalidArgs, " "))
	}

	config := Config{
		Reverse:      strings.ContainsRune(flatArgs, FlagReverse),
		CandiedYAML:  strings.ContainsRune(flatArgs, FlagCandiedYAML),
		JSONDecoder:  strings.ContainsRune(flatArgs, FlagJSONDecoder),
		EscapeHTML:   strings.ContainsRune(flatArgs, FlagEscapeHTML),
		FloatStrings: !strings.ContainsRune(flatArgs, FlagNoFloatStrings),
		JSONKeys:     strings.ContainsRune(flatArgs, FlagJSONKeys),
		Help:         strings.ContainsRune(flatArgs, FlagHelp),
	}

	if !config.Reverse && config.JSONDecoder {
		return Config{}, fmt.Errorf("flag -%c cannot be specified without flag -%c", FlagJSONDecoder, FlagReverse)
	}

	if !config.Reverse && config.JSONKeys {
		return Config{}, fmt.Errorf("flag -%c cannot be specified without flag -%c", FlagJSONKeys, FlagReverse)
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagReverse, FlagCandiedYAML, FlagJSONDecoder, FlagEscapeHTML,
		FlagNoFloatStrings, FlagJSONKeys, FlagHelp, '\t', ' ', '-':
		return -1
	}
	return r
}
