package yaml_test

import (
	"testing"

	"github.com/sclevine/yj/yaml"
)

func TestJSONMarshal(t *testing.T) {
	json := &yaml.JSON{}
	out, err := json.Marshal(F{`some-<\u003c`})
	assertEq(t, err, nil)
	assertEq(t, string(out), `{"f":"some-<\\u003c"}`)

	json.EscapeHTML = true
	out, err = json.Marshal(F{`some-<\u003c`})
	assertEq(t, err, nil)
	assertEq(t, string(out), `{"f":"some-\u003c\\u003c"}`)
}

func TestJSONMarshalWhenInputIsInvalid(t *testing.T) {
	json := &yaml.JSON{}
	_, err := json.Marshal(func() {})
	assertEq(t, err.Error(), "json: unsupported type: func()")
}
