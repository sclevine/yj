package yaml_test

import (
	"errors"
	"testing"

	"github.com/sclevine/yj/yaml"
)

func TestDecoderPanics(t *testing.T) {
	var panicValue interface{}
	decoder := &yaml.Decoder{DecodeYAML: func(_ interface{}) error {
		panic(panicValue)
	}}

	panicValue = errors.New("some error")
	_, err := decoder.JSON()
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	_, err = decoder.JSON()
	assertEq(t, err.Error(), "unexpected failure: some panic")
}

func TestEncoderPanics(t *testing.T) {
	var panicValue interface{}
	encoder := &yaml.Encoder{EncodeYAML: func(_ interface{}) error {
		panic(panicValue)
	}}

	panicValue = errors.New("some error")
	err := encoder.YAML(nil)
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	err = encoder.YAML(nil)
	assertEq(t, err.Error(), "unexpected failure: some panic")
}
