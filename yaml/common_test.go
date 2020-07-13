package yaml_test

import (
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type mockYAML struct {
	data  []byte
	value interface{}
	err   error
}

func (m *mockYAML) decode(r io.Reader) (interface{}, error) {
	var err error
	m.data, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return m.value, m.err
}

func (m *mockYAML) encode(w io.Writer, v interface{}) error {
	if _, err := w.Write(m.data); err != nil {
		return err
	}
	m.value = v
	return m.err
}

func assertEq(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
