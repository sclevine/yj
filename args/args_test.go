package args_test

import (
	"reflect"
	"testing"

	"github.com/sclevine/yj/args"
	"github.com/sclevine/yj/convert"
)

func TestParse(t *testing.T) {
	config, err := args.Parse("-t", "y\tk-", "kn-k ", "h h", "")
	assertEq(t, err, nil)
	_, ok := config.From.(convert.TOML)
	assertEq(t, ok, true)
	yaml, ok := config.To.(convert.YAML)
	assertEq(t, ok, true)
	assertEq(t, yaml, convert.YAML{
		EscapeHTML:   false,
		FloatStrings: false,
		JSONKeys:     true,
	})
	assertEq(t, config.Help, true)

	config, err = args.Parse("--\t\te  ", "")
	assertEq(t, err, nil)
	yaml, ok = config.From.(convert.YAML)
	assertEq(t, ok, true)
	assertEq(t, yaml, convert.YAML{
		EscapeHTML:   true,
		FloatStrings: true,
		JSONKeys:     false,
	})
	json, ok := config.To.(convert.JSON)
	assertEq(t, ok, true)
	assertEq(t, json, convert.JSON{
		EscapeHTML: true,
	})
	assertEq(t, config.Help, false)

	// TODO: test more ytj combinations
}

func TestParseWithInvalidFlags(t *testing.T) {
	_, err := args.Parse("-ar", "yb\te-", "kn-k ", "h ~ dh", "ab")
	assertEq(t, err.Error(), "invalid flags specified: a b ~ d a b")

	_, err = args.Parse("k")
	assertEq(t, err.Error(), "flag -k only valid for YAML output")

	_, err = args.Parse("ejy")
	assertEq(t, err.Error(), "flag -e only valid for JSON output")
}

func assertEq(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
