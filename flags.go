package main

import (
	"fmt"
	"strings"

	"github.com/sclevine/yj/convert"
)

type Config struct {
	From, To      convert.Encoding
	Help, Version bool
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
	FlagVersion        = 'v'
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
		From:    from,
		To:      to,
		Help:    strings.ContainsRune(flatArgs, FlagHelp),
		Version: strings.ContainsRune(flatArgs, FlagVersion),
	}

	return config, nil
}

func flagFilter(r rune) rune {
	switch r {
	case FlagYAML, FlagTOML, FlagJSON, FlagHCL, FlagReverse,
		FlagEscapeHTML, FlagNoFloatStrings, FlagJSONKeys,
		FlagIndent, FlagHelp, FlagVersion, '\t', ' ', '-':
		return -1
	}
	return r
}

func transform(s string) (from, to convert.Encoding, err error) {
	escapeHTML := strings.ContainsRune(s, FlagEscapeHTML)
	indent := strings.ContainsRune(s, FlagIndent)
	jsonKeys := strings.ContainsRune(s, FlagJSONKeys)
	floatStrings := !strings.ContainsRune(s, FlagNoFloatStrings)

	yaml := &convert.YAML{
		JSONKeys:   jsonKeys,
		EscapeHTML: escapeHTML,
	}
	toml := &convert.TOML{
		Indent: indent,
	}
	json := &convert.JSON{
		EscapeHTML: escapeHTML,
		Indent:     indent,
	}
	hcl := &convert.HCL{}

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

	setKeyFloats(to, func(setTo func(v convert.SpecialFloats)) {
		if floatStrings {
			setTo(convert.FloatsString)
		} else {
			// never convert number to nan, inf, -inf
			setTo(convert.FloatsReal)
		}
	})

	setKeyFloats(from, func(setFrom func(v convert.SpecialFloats)) {
		if floatStrings {
			setFrom(convert.FloatsString)
		} else {
			setFrom(convert.FloatsNumber)
		}
	})

	setFloats(to, func(setTo func(v convert.SpecialFloats)) {
		if floatStrings {
			setTo(convert.FloatsString)
		} else {
			// never convert number to nan, inf, -inf
			setTo(convert.FloatsReal)
		}
	})

	setFloats(from, func(setFrom func(v convert.SpecialFloats)) {
		if floatStrings {
			setFrom(convert.FloatsString)
		} else {
			setFrom(convert.FloatsNumber)
		}
		setFloats(to, func(setTo func(v convert.SpecialFloats)) {
			setFrom(convert.FloatsReal)
			setTo(convert.FloatsReal)
		})
	})

	if _, toYAML := to.(*convert.YAML); jsonKeys && !toYAML {
		err = fmt.Errorf("flag -%c only valid for YAML output", FlagJSONKeys)
		return
	}
	if _, toJSON := to.(*convert.JSON); escapeHTML && !toJSON {
		err = fmt.Errorf("flag -%c only valid for JSON output", FlagEscapeHTML)
		return
	}

	if indent {
		switch to.(type) {
		case *convert.JSON, *convert.TOML:
		default:
			err = fmt.Errorf("flag -%c only valid for JSON or TOML output", FlagIndent)
			return
		}
	}

	// TODO: validate -n isn't used between inapplicable types

	return
}

func setFloats(e convert.Encoding, f func(set func(v convert.SpecialFloats))) {
	switch e := e.(type) {
	case *convert.YAML:
		f(func(v convert.SpecialFloats) { e.SpecialFloats = v })
	case *convert.TOML:
		f(func(v convert.SpecialFloats) { e.SpecialFloats = v })
	}
}

func setKeyFloats(e convert.Encoding, f func(set func(v convert.SpecialFloats))) {
	switch e := e.(type) {
	case *convert.YAML:
		f(func(v convert.SpecialFloats) { e.KeySpecialFloats = v })
	}
}
