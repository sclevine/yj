package yaml_test

import (
	"bytes"
	"errors"
	"testing"

	goyaml "gopkg.in/yaml.v2"

	"github.com/sclevine/yj/yaml"
)

func TestEncoder(t *testing.T) {
	mock := &mockYAML{}
	encoder := &yaml.Encoder{
		EncodeYAML:   mock.encode,
		KeyUnmarshal: keyUnmarshal,
		NaN:          F{"NaN"},
		PosInf:       F{"Infinity"},
		NegInf:       F{"-Infinity"},
	}
	err := encoder.YAML(jsonFixture)
	assertEq(t, err, nil)
	assertEq(t, mock.value, yamlFixture)
}

func TestEncoderWhenYAMLIsInvalid(t *testing.T) {
	mock := &mockYAML{err: errors.New("some error")}
	encoder := &yaml.Encoder{EncodeYAML: mock.encode}
	err := encoder.YAML(nil)
	assertEq(t, err.Error(), "some error")
}

func TestEncoderWhenYAMLHasInvalidTypes(t *testing.T) {
	mock := &mockYAML{}
	encoder := &yaml.Encoder{EncodeYAML: mock.encode}

	err := encoder.YAML(map[int]int{})
	assertEq(t, err.Error(), "unexpected type: map[int]int{}")

	err = encoder.YAML([0]int{})
	assertEq(t, err.Error(), "unexpected type: [0]int{}")

	err = encoder.YAML([]int{})
	assertEq(t, err.Error(), "unexpected type: []int{}")

	err = encoder.YAML(float32(0))
	assertEq(t, err.Error(), "unexpected type: 0")
}

func keyUnmarshal(data []byte, v interface{}) error {
	switch {
	case bytes.HasPrefix(data, []byte("stringer=")):
		*v.(*interface{}) = stringer(string(data[9:]))
	case bytes.HasPrefix(data, []byte(`{"f":`)):
		key := &F{}
		if err := goyaml.Unmarshal(data, key); err != nil {
			return err
		}
		*v.(*interface{}) = *key
	default:
		var key interface{}
		if err := goyaml.Unmarshal(data, &key); err != nil {
			return err
		}
		*v.(*interface{}) = key
	}
	return nil
}

// Test pointer keys for map/slice/func
// Test nil values not converted to NaN
// Test without key unmarshaling
