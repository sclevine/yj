package yaml_test

import (
	"errors"
	"testing"

	"github.com/sclevine/yj/yaml"
)

func TestDecoderPanics(t *testing.T) {
	var panicValue interface{}
	decoder := &yaml.Decoder{Unmarshal: func(_ []byte, _ interface{}) error {
		panic(panicValue)
	}}

	panicValue = errors.New("some error")
	_, err := decoder.JSON(nil)
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	_, err = decoder.JSON(nil)
	assertEq(t, err.Error(), "unexpected failure: some panic")
}

func TestEncoderPanics(t *testing.T) {
	var panicValue interface{}
	encoder := &yaml.Encoder{Marshal: func(_ interface{}) ([]byte, error) {
		panic(panicValue)
	}}

	panicValue = errors.New("some error")
	_, err := encoder.YAML(nil)
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	_, err = encoder.YAML(nil)
	assertEq(t, err.Error(), "unexpected failure: some panic")
}
