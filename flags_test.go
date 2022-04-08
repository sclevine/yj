package main_test

import (
	"testing"

	"github.com/sclevine/yj/v5"
	"github.com/sclevine/yj/v5/convert"
)

func TestParse(t *testing.T) {
	config, err := main.Parse("-t", "y\tk-", "kn-k ", "h h", "")
	assertEq(t, err, nil)
	toml, ok := config.From.(*convert.TOML)
	assertEq(t, ok, true)
	assertEq(t, toml, &convert.TOML{
		SpecialFloats: convert.FloatsReal,
	})
	yaml, ok := config.To.(*convert.YAML)
	assertEq(t, ok, true)
	assertEq(t, yaml, &convert.YAML{
		SpecialFloats: convert.FloatsReal,
		EscapeHTML:    false,
		JSONKeys:      true,
	})
	assertEq(t, config.Help, true)

	config, err = main.Parse("--\t\te  ", "")
	assertEq(t, err, nil)
	yaml, ok = config.From.(*convert.YAML)
	assertEq(t, ok, true)
	assertEq(t, yaml, &convert.YAML{
		SpecialFloats:    convert.FloatsString,
		KeySpecialFloats: convert.FloatsString,
		EscapeHTML:       true,
		JSONKeys:         false,
	})
	json, ok := config.To.(*convert.JSON)
	assertEq(t, ok, true)
	assertEq(t, json, &convert.JSON{
		EscapeHTML: true,
	})
	assertEq(t, config.Help, false)

	// TODO: test more ytjc combinations
}

func TestParseWithInvalidFlags(t *testing.T) {
	_, err := main.Parse("-ar", "yb\te-", "kn-k ", "h ~ dh", "ab")
	assertEq(t, err.Error(), "invalid flags specified: a b ~ d a b")

	_, err = main.Parse("k")
	assertEq(t, err.Error(), "flag -k only valid for YAML output")

	_, err = main.Parse("ejy")
	assertEq(t, err.Error(), "flag -e only valid for JSON output")

	_, err = main.Parse("ijy")
	assertEq(t, err.Error(), "flag -i only valid for JSON or TOML output")
}
