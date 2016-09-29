package yaml_test

import (
	"bytes"
	"errors"
	"testing"

	goyaml "gopkg.in/yaml.v2"

	"github.com/sclevine/yj/yaml"
)

func TestEncoderWhenYAMLIsInvalid(t *testing.T) {
	mock := &mockYAML{err: errors.New("some error")}
	encoder := &yaml.Encoder{Marshal: mock.marshal}
	_, err := encoder.YAML(nil)
	assertEqual(t, err.Error(), "some error")
}

func TestEncoderWhenYAMLHasInvalidTypes(t *testing.T) {
	mock := &mockYAML{}
	encoder := &yaml.Encoder{Marshal: mock.marshal}

	_, err := encoder.YAML(map[int]int{})
	assertEqual(t, err.Error(), "unexpected type: map[int]int{}")

	_, err = encoder.YAML([0]int{})
	assertEqual(t, err.Error(), "unexpected type: [0]int{}")

	_, err = encoder.YAML([]int{})
	assertEqual(t, err.Error(), "unexpected type: []int{}")

	_, err = encoder.YAML(float32(0))
	assertEqual(t, err.Error(), "unexpected type: 0")
}

func TestEncoderWhenYAMLIsValid(t *testing.T) {
	mock := &mockYAML{data: []byte("some YAML")}
	encoder := &yaml.Encoder{
		Marshal:      mock.marshal,
		KeyUnmarshal: keyUnmarshal,
		NaN:          F{"NaN"},
		PosInf:       F{"Infinity"},
		NegInf:       F{"-Infinity"},
	}
	yaml, err := encoder.YAML(jsonFixture)
	assertEqual(t, err, nil)
	assertEqual(t, mock.value, yamlFixture)
	assertEqual(t, yaml, []byte("some YAML"))
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
