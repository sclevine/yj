package args

import (
	"fmt"
	"strings"
)

type Config struct {
	Reverse      bool
	CandiedYAML  bool
	EscapeHTML   bool
	FloatStrings bool
	JSONKeys     bool
}

const (
	FlagReverse        = 'r'
	FlagCandiedYAML    = 'c'
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
		EscapeHTML:   strings.ContainsRune(flatArgs, FlagEscapeHTML),
		FloatStrings: !strings.ContainsRune(flatArgs, FlagNoFloatStrings),
		JSONKeys:     strings.ContainsRune(flatArgs, FlagJSONKeys),
	}

	if !config.Reverse && config.JSONKeys {
		return Config{}, fmt.Errorf("flag -%c cannot be specified without flag -%c", FlagJSONKeys, FlagReverse)
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagReverse, FlagCandiedYAML, FlagEscapeHTML,
		FlagNoFloatStrings, FlagJSONKeys, FlagHelp, '\t', ' ', '-':
		return -1
	}
	return r
}
