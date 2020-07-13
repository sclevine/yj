package yaml_test

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/sclevine/yj/yaml"
)

func TestDecoder(t *testing.T) {
	mock := &mockYAML{value: yamlFixture}
	decoder := &yaml.Decoder{
		DecodeYAML: mock.decode,
		KeyMarshal: json.Marshal,
		NaN:        F{"NaN"},
		PosInf:     F{"Infinity"},
		NegInf:     F{"-Infinity"},
	}
	r := strings.NewReader("test")
	json, err := decoder.JSON(r)
	assertEq(t, err, nil)
	assertEq(t, json, jsonFixture)
	assertEq(t, mock.data, []byte("test"))
}

func TestDecoderWhenYAMLIsInvalid(t *testing.T) {
	mock := &mockYAML{err: errors.New("some error")}
	decoder := &yaml.Decoder{DecodeYAML: mock.decode}
	r := strings.NewReader("test")
	_, err := decoder.JSON(r)
	assertEq(t, err.Error(), "some error")
}

func TestDecoderWhenYAMLHasInvalidTypes(t *testing.T) {
	mock := &mockYAML{}
	decoder := &yaml.Decoder{DecodeYAML: mock.decode}
	r := strings.NewReader("test")

	mock.value = map[int]int{}
	_, err := decoder.JSON(r)
	assertEq(t, err.Error(), "unexpected type: map[int]int{}")

	mock.value = [0]int{}
	_, err = decoder.JSON(r)
	assertEq(t, err.Error(), "unexpected type: [0]int{}")

	mock.value = []int{}
	_, err = decoder.JSON(r)
	assertEq(t, err.Error(), "unexpected type: []int{}")

	mock.value = float32(0)
	_, err = decoder.JSON(r)
	assertEq(t, err.Error(), "unexpected type: 0")
}

func TestDecoderWhenYAMLHasInvalidKeys(t *testing.T) {
	mock := &mockYAML{value: map[interface{}]interface{}{
		math.NaN(): "",
	}}
	decoder := &yaml.Decoder{
		DecodeYAML: mock.decode,
		KeyMarshal: json.Marshal,
	}
	r := strings.NewReader("test")
	_, err := decoder.JSON(r)
	assertEq(t, err.Error(), "json: unsupported value: NaN")
}
