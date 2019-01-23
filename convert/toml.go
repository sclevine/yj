package convert

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type TOML struct{}

func (TOML) String() string {
	return "TOML"
}

func (TOML) Encode(input interface{}) ([]byte, error) {
	output := &bytes.Buffer{}
	err := toml.NewEncoder(output).Encode(input)
	return output.Bytes(), err
}

func (TOML) Decode(input []byte) (interface{}, error) {
	var data interface{}
	return data, toml.Unmarshal(input, &data)
}
