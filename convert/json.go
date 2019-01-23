package convert

import (
	"bytes"
	"encoding/json"
)

type JSON struct {
	EscapeHTML bool
}

func (JSON) String() string {
	return "JSON"
}

func (j JSON) Encode(input interface{}) ([]byte, error) {
	output := &bytes.Buffer{}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(j.EscapeHTML)
	err := encoder.Encode(input)
	return output.Bytes(), err
}

func (JSON) Decode(input []byte) (interface{}, error) {
	var data interface{}
	return data, json.Unmarshal(input, &data)
}
