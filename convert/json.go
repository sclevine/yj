package convert

import (
	"encoding/json"
	"io"
)

type JSON struct {
	EscapeHTML bool
	Indent     bool
}

func (JSON) String() string {
	return "JSON"
}

func (j JSON) Encode(w io.Writer, in interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(j.EscapeHTML)
	if j.Indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(in)
}

func (JSON) Decode(r io.Reader) (interface{}, error) {
	var data interface{}
	return data, json.NewDecoder(r).Decode(&data)
}
