package main_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	main "github.com/sclevine/yj/v5"
)

func TestRunWhenArgsFailToParse(t *testing.T) {
	stdout, stderr, code := run(nil, "some-bad-args")
	assertEq(t, code, 1)
	assertEq(t, bytes.Contains(stderr, []byte("Error: invalid flags")), true)
	assertEq(t, bytes.Contains(stderr, []byte("Usage:")), true)
	assertEq(t, len(stdout), 0)
}

func TestRunWhenHelpFlagIsProvided(t *testing.T) {
	stdout, stderr, code := run(nil, "hr")
	assertEq(t, code, 0)
	assertEq(t, bytes.Contains(stdout, []byte("Usage:")), true)
	assertEq(t, len(stderr), 0)
}

func TestRunWhenStdinIsInvalid(t *testing.T) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	assertEq(t, main.Run(&errReader{}, stdout, stderr, []string{""}), 1)
	assertEq(t, strings.Contains(stderr.String(), "some reader error"), true)
	assertEq(t, stdout.Len(), 0)
}

func TestCases(t *testing.T) {
	outs, err := filepath.Glob("testdata/case*_out_*")
	assertEq(t, err, nil)
	for _, out := range outs {
		t.Log(out)
		parts := strings.SplitN(out, "_out_", 2)
		flags := strings.SplitN(parts[1], ".", 2)
		ins, err := filepath.Glob(parts[0] + "_in.*")
		assertEq(t, err, nil)
		stdout, stderr, code := run(rdfile(t, ins[0]), flags[0])
		assertEq(t, string(stderr), "")
		assertEq(t, code, 0)
		assertEq(t, string(rdfile(t, out)), string(stdout))
	}
}

func TestTOMLLiteralKeys(t *testing.T) {
	cases := []struct {
		name, input string
	}{
		{"comma", `{"stack(x,n)":"foobar"}`},
		{"empty and minus", `{"":1,"-":2,",":3,"a,":4}`},
		{"shared prefix", `{"a":1,"a,b":2,"a,c":3}`},
		{"tag options", `{"a,omitempty":0,"a,omitzero":false}`},
		{"nested", `{"outer,key":{"inner,key":1,"":2,"-":3}}`},
		{"table array", `{"array,key":[{"inner,key":1},{"inner,key":2,"":3}]}`},
		{"inline table", `{"array,key":[1,{"inner,key":2}]}`},
		{"escaped", `{"quote\",key":1,"slash\\,key":2,"line\n,key":3,"tab\t,key":4,"雪,key":5}`},
		{"alias text", `{"a,b":"a,b","_yj_key_0":{"_yj_key_1":"_yj_key_0"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, flags := range []string{"jt", "jti"} {
				t.Run(flags, func(t *testing.T) {
					encoded, stderr, code := run([]byte(tc.input), flags)
					assertEq(t, string(stderr), "")
					assertEq(t, code, 0)
					decoded, stderr, code := run(encoded, "tj")
					assertEq(t, string(stderr), "")
					assertEq(t, code, 0)
					assertEq(t, string(decoded), tc.input+"\n")
				})
			}
		})
	}
}

func run(in []byte, flags ...string) (stdout, stderr []byte, code int) {
	out, err := &bytes.Buffer{}, &bytes.Buffer{}
	code = main.Run(bytes.NewReader(in), out, err, append([]string{"yj"}, flags...))
	return out.Bytes(), err.Bytes(), code
}

func rdfile(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	assertEq(t, err, nil)
	return b
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some reader error")
}

func assertEq(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
