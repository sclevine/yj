package yaml

import (
	"bytes"
	"encoding/json"
)

type JSON struct {
	EscapeHTML bool
}

func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	keyJSON := &bytes.Buffer{}
	encoder := json.NewEncoder(keyJSON)
	encoder.SetEscapeHTML(j.EscapeHTML)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return keyJSON.Bytes()[:keyJSON.Len()-1], nil
}
