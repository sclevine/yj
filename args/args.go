package args

import (
	"fmt"
	"strings"
)

type Config struct {
	Reverse      bool
	CandiedYAML  bool
	JSONAsYAML   bool
	EscapeHTML   bool
	FloatStrings bool
	JSONKeys     bool
	Help         bool
}

const (
	FlagReverse        = 'r'
	FlagCandiedYAML    = 'c'
	FlagJSONAsYAML     = 'y'
	FlagEscapeHTML     = 'e'
	FlagNoFloatStrings = 'n'
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
		Reverse:      strings.ContainsRune(flatArgs, FlagReverse),
		CandiedYAML:  strings.ContainsRune(flatArgs, FlagCandiedYAML),
		JSONAsYAML:   strings.ContainsRune(flatArgs, FlagJSONAsYAML),
		EscapeHTML:   strings.ContainsRune(flatArgs, FlagEscapeHTML),
		FloatStrings: !strings.ContainsRune(flatArgs, FlagNoFloatStrings),
		JSONKeys:     strings.ContainsRune(flatArgs, FlagJSONKeys),
		Help:         strings.ContainsRune(flatArgs, FlagHelp),
	}

	if !config.Reverse && config.JSONAsYAML {
		return nil, fmt.Errorf("flag -%c cannot be specified without flag -%c", FlagJSONAsYAML, FlagReverse)
	}

	if !config.Reverse && config.JSONKeys {
		return nil, fmt.Errorf("flag -%c cannot be specified without flag -%c", FlagJSONKeys, FlagReverse)
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagReverse, FlagCandiedYAML, FlagJSONAsYAML, FlagEscapeHTML,
		FlagNoFloatStrings, FlagJSONKeys, FlagHelp, '\t', ' ', '-':
		return -1
	}
	return r
}
