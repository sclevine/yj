package main_test

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/sclevine/yj"
)

func TestRunWhenArgsFailToParse(t *testing.T) {
	stdin, stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	assertEqual(t, main.Run(stdin, stdout, stderr, []string{"", "some-bad-args"}), 1)
	assertEqual(t, strings.Contains(stderr.String(), "Error: invalid flags"), true)
	assertEqual(t, strings.Contains(stderr.String(), "Usage:"), true)
	assertEqual(t, stdout.Len(), 0)
}

func TestRunWhenHelpFlagIsProvided(t *testing.T) {
	stdin, stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	assertEqual(t, main.Run(stdin, stdout, stderr, []string{"", "hr"}), 0)
	assertEqual(t, strings.Contains(stdout.String(), "Usage:"), true)
	assertEqual(t, stderr.Len(), 0)
}

func TestRunWhenStdinIsInvalid(t *testing.T) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	assertEqual(t, main.Run(&errReader{}, stdout, stderr, []string{""}), 1)
	assertEqual(t, strings.Contains(stderr.String(), "Error: some reader error"), true)
	assertEqual(t, strings.Contains(stderr.String(), "Usage:"), true)
	assertEqual(t, stdout.Len(), 0)
}

// Needs integration tests

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some reader error")
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("\nAssertion failed:\n\t%#v\nnot equal to\n\t%#v\n", a, b)
	}
}
