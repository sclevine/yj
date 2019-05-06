package main

import (
	"fmt"
	"strings"

	"github.com/sclevine/yj/convert"
)

type Config struct {
	From, To convert.Encoding
	Help     bool
}

const (
	FlagYAML           = 'y'
	FlagTOML           = 't'
	FlagJSON           = 'j'
	FlagHCL            = 'c'
	FlagReverse        = 'r'
	FlagNoFloatStrings = 'n'
	FlagEscapeHTML     = 'e'
	FlagIndent         = 'i'
	FlagJSONKeys       = 'k'
	FlagHelp           = 'h'
)

func Parse(args ...string) (*Config, error) {
	flatArgs := strings.Join(args, "")
	invalidArgs := strings.Split(strings.Map(flagFilter, flatArgs), "")

	if len(invalidArgs) > 0 {
		return nil, fmt.Errorf("invalid flags specified: %s", strings.Join(invalidArgs, " "))
	}

	from, to, err := transform(flatArgs)
	if err != nil {
		return nil, err
	}
	config := &Config{
		From: from,
		To:   to,
		Help: strings.ContainsRune(flatArgs, FlagHelp),
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagYAML, FlagTOML, FlagJSON, FlagHCL, FlagReverse,
		FlagEscapeHTML, FlagNoFloatStrings, FlagJSONKeys,
		FlagIndent, FlagHelp, '\t', ' ', '-':
		return -1
	}
	return r
}

func transform(s string) (from, to convert.Encoding, err error) {
	escapeHTML := strings.ContainsRune(s, FlagEscapeHTML)
	indent := strings.ContainsRune(s, FlagIndent)
	floatStrings := !strings.ContainsRune(s, FlagNoFloatStrings)
	jsonKeys := strings.ContainsRune(s, FlagJSONKeys)

	yaml := convert.YAML{
		FloatStrings: floatStrings,
		JSONKeys:     jsonKeys,
		EscapeHTML:   escapeHTML,
	}
	toml := convert.TOML{
		Indent: indent,
	}
	json := convert.JSON{
		EscapeHTML: escapeHTML,
		Indent:     indent,
	}
	hcl := convert.HCL{}

	for _, r := range s {
		switch r {
		case FlagYAML:
			from, to = to, yaml
		case FlagTOML:
			from, to = to, toml
		case FlagJSON:
			from, to = to, json
		case FlagHCL:
			from, to = to, hcl
		case FlagReverse:
			from, to = json, yaml
		}
	}
	if from == nil {
		if to == nil {
			to = yaml
		}
		from, to = to, json
	}

	if _, toYAML := to.(convert.YAML); jsonKeys && !toYAML {
		err = fmt.Errorf("flag -%c only valid for YAML output", FlagJSONKeys)
	}
	if _, toJSON := to.(convert.JSON); escapeHTML && !toJSON {
		err = fmt.Errorf("flag -%c only valid for JSON output", FlagEscapeHTML)
	}

	switch to.(type) {
	case convert.JSON, convert.TOML:
	default:
		if indent {
			err = fmt.Errorf("flag -%c only valid for JSON or TOML output", FlagIndent)
		}
	}

	// TODO: validate -n has YAML input or output flag

	return
}
