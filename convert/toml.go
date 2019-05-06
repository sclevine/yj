package convert

import (
	"io"

	"github.com/BurntSushi/toml"
)

type TOML struct {
	Indent bool
}

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) error {
	enc := toml.NewEncoder(w)
	if !t.Indent {
		enc.Indent = ""
	}
	return enc.Encode(in)
}

func (TOML) Decode(r io.Reader) (interface{}, error) {
	var data interface{}
	_, err := toml.DecodeReader(r, &data)
	return data, err
}
