package args_test

import (
	"reflect"
	"testing"

	"github.com/sclevine/yj/args"
)

func TestParse(t *testing.T) {
	config, err := args.Parse("-rc", "y\te-", "kn-k ", "h h", "")
	assertEqual(t, err, nil)
	assertEqual(t, config, &args.Config{
		Reverse:      true,
		CandiedYAML:  true,
		JSONAsYAML:   true,
		EscapeHTML:   true,
		FloatStrings: false,
		JSONKeys:     true,
		Help:         true,
	})
	config, err = args.Parse("--\t\t  ", "")
	assertEqual(t, err, nil)
	assertEqual(t, config, &args.Config{
		Reverse:      false,
		CandiedYAML:  false,
		JSONAsYAML:   false,
		EscapeHTML:   false,
		FloatStrings: true,
		JSONKeys:     false,
		Help:         false,
	})
}

func TestParseWithInvalidFlags(t *testing.T) {
	_, err := args.Parse("-arc", "yb\te-", "kn-k ", "hc ~ dh", "ab")
	assertEqual(t, err.Error(), "invalid flags specified: a b ~ d a b")

	_, err = args.Parse("y")
	assertEqual(t, err.Error(), "flag -y cannot be specified without flag -r")

	_, err = args.Parse("k")
	assertEqual(t, err.Error(), "flag -k cannot be specified without flag -r")
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
