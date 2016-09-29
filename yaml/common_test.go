package yaml_test

import (
	"reflect"
	"testing"
)

type mockYAML struct {
	data  []byte
	value interface{}
	err   error
}

func (m *mockYAML) unmarshal(data []byte, v interface{}) error {
	m.data = data
	*v.(*interface{}) = m.value
	return m.err
}

func (m *mockYAML) marshal(v interface{}) ([]byte, error) {
	m.value = v
	return m.data, m.err
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
