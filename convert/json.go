package convert

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/sclevine/yj/order"
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
	// TODO: remove global, may affect yaml/json.go
	defer func(p bool) {
		order.MapSliceEscapeHTML = p
	}(order.MapSliceEscapeHTML)
	order.MapSliceEscapeHTML = j.EscapeHTML
	if j.Indent {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(in)
}

// TODO: implement streaming version
func (JSON) Decode(r io.Reader) (interface{}, error) {
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if !json.Valid(out) {
		var null interface{}
		err := json.Unmarshal(out, &null)
		if err == nil {
			err = errors.New("invalid JSON")
		}
		return nil, err
	}
	return (YAML{}).Decode(bytes.NewReader(out))
}
