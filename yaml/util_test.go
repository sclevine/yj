package yaml_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/sclevine/yj/yaml"
)

func TestDecoderPanics(t *testing.T) {
	var panicValue interface{}
	decoder := &yaml.Decoder{DecodeYAML: func(_ io.Reader) (interface{}, error) {
		panic(panicValue)
	}}
	r := strings.NewReader("test")

	panicValue = errors.New("some error")
	_, err := decoder.JSON(r)
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	_, err = decoder.JSON(r)
	assertEq(t, err.Error(), "unexpected failure: some panic")
}

func TestEncoderPanics(t *testing.T) {
	var panicValue interface{}
	encoder := &yaml.Encoder{EncodeYAML: func(_ io.Writer, _ interface{}) error {
		panic(panicValue)
	}}
	w := &bytes.Buffer{}

	panicValue = errors.New("some error")
	err := encoder.YAML(w, nil)
	assertEq(t, err.Error(), "some error")

	panicValue = "some panic"
	err = encoder.YAML(w, nil)
	assertEq(t, err.Error(), "unexpected failure: some panic")
}
