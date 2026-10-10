package convert

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/sclevine/yj/v5/order"
)

func TestRestoreTOMLKeys(t *testing.T) {
	input := "[_key0]\n_key1 = \"_key0 \\\"_key1\\\"\"\n" +
		"_key0_suffix = ['_key1', \"\"\"_key0\n_key1\"\"\", '''_key0\n_key1''']\n" +
		"_key2 = [\"\"\"_key0\"\"\"\", \"\"\"_key1\"\"\"\"\", '''_key0'''', '''_key1''''']\n" +
		"# _key0 = _key1\n_key2 = [{_key1 = \"_key0\"}]\n"
	want := "[\"a,b\"]\n_key0 = \"_key0 \\\"_key1\\\"\"\n" +
		"_key0_suffix = ['_key1', \"\"\"_key0\n_key1\"\"\", '''_key0\n_key1''']\n" +
		"\"\" = [\"\"\"_key0\"\"\"\", \"\"\"_key1\"\"\"\"\", '''_key0'''', '''_key1''''']\n" +
		"# _key0 = _key1\n\"\" = [{_key0 = \"_key0\"}]\n"
	var out bytes.Buffer
	if err := restoreTOMLKeys(&out, []byte(input), map[string]string{
		"_key0": `"a,b"`, "_key1": "_key0", "_key2": `""`,
	}); err != nil {
		t.Fatal(err)
	}
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func TestTOMLAliasCollisions(t *testing.T) {
	alias := tomlKeyPrefix + "0"
	in := order.MapSlice{{Key: alias, Val: alias}, {Key: "a,b", Val: "a,b"}}
	var out bytes.Buffer
	if err := (TOML{}).Encode(&out, in); err != nil {
		t.Fatal(err)
	}
	want := alias + " = \"" + alias + "\"\n\"a,b\" = \"a,b\"\n"
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
	if !reflect.DeepEqual(in, order.MapSlice{{Key: alias, Val: alias}, {Key: "a,b", Val: "a,b"}}) {
		t.Fatal("encoding changed the input keys or values")
	}
}

func TestTOMLRepeatedEncoding(t *testing.T) {
	in := order.MapSlice{{Key: "a,b", Val: "value"}}
	for i := 0; i < 2; i++ {
		var out bytes.Buffer
		if err := (TOML{}).Encode(&out, in); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(out.String(), tomlKeyPrefix) {
			t.Fatalf("alias remains in output %q", out.String())
		}
		if out.String() != "\"a,b\" = \"value\"\n" {
			t.Fatalf("unexpected output %q", out.String())
		}
	}
}

func TestTOMLWriterErrors(t *testing.T) {
	for _, err := range []error{nil, errors.New("write failed")} {
		writer := &tomlErrorWriter{err: err}
		got := (TOML{}).Encode(writer, order.MapSlice{{Key: "key", Val: "value"}})
		want := err
		if want == nil {
			want = io.ErrShortWrite
		}
		if !errors.Is(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

type tomlErrorWriter struct {
	err error
}

func (w *tomlErrorWriter) Write(p []byte) (int, error) {
	return 0, w.err
}
