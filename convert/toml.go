package convert

import (
	"io"

	"github.com/BurntSushi/toml"
)

type TOML struct{}

func (TOML) String() string {
	return "TOML"
}

func (TOML) Encode(w io.Writer, in interface{}) error {
	return toml.NewEncoder(w).Encode(in)
}

func (TOML) Decode(r io.Reader) (interface{}, error) {
	var data interface{}
	_, err := toml.DecodeReader(r, &data)
	return data, err
}
