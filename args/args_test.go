package args_test

import (
	"reflect"
	"testing"

	"github.com/sclevine/yj/args"
)

func TestParse(t *testing.T) {
	config, err := args.Parse("-t", "y\tk-", "kn-k ", "h h", "")
	assertEqual(t, err, nil)
	assertEqual(t, config, &args.Config{
		From:         args.TOML,
		To:           args.YAML,
		EscapeHTML:   false,
		FloatStrings: false,
		JSONKeys:     true,
		Help:         true,
	})
	config, err = args.Parse("--\t\te  ", "")
	assertEqual(t, err, nil)
	assertEqual(t, config, &args.Config{
		From:         args.YAML,
		To:           args.JSON,
		EscapeHTML:   true,
		FloatStrings: true,
		JSONKeys:     false,
		Help:         false,
	})
	// TODO: test more ytj combinations
}

func TestParseWithInvalidFlags(t *testing.T) {
	_, err := args.Parse("-ar", "yb\te-", "kn-k ", "h ~ dh", "ab")
	assertEqual(t, err.Error(), "invalid flags specified: a b ~ d a b")

	_, err = args.Parse("k")
	assertEqual(t, err.Error(), "flag -k only valid for YAML output")

	_, err = args.Parse("ejy")
	assertEqual(t, err.Error(), "flag -e only valid for JSON output")
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
