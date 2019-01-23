package yaml_test

import (
	"reflect"
	"testing"
)

type mockYAML struct {
	value interface{}
	err   error
}

func (m *mockYAML) decode(v interface{}) error {
	*v.(*interface{}) = m.value
	return m.err
}

func (m *mockYAML) encode(v interface{}) error {
	m.value = v
	return m.err
}

func assertEq(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
