package yaml

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"

	goyaml "gopkg.in/yaml.v3"
)

type KeyJSON struct {
	EscapeHTML bool
}

func (k *KeyJSON) Marshal(v interface{}) ([]byte, error) {
	keyJSON := &bytes.Buffer{}
	jsonEnc := json.NewEncoder(keyJSON)
	jsonEnc.SetEscapeHTML(k.EscapeHTML)
	if err := jsonEnc.Encode(v); err != nil {
		return nil, err
	}
	return keyJSON.Bytes()[:keyJSON.Len()-1], nil
}

func (k *KeyJSON) Unmarshal(src []byte, v interface{}) error {
	if !json.Valid(src) {
		var null interface{}
		err := json.Unmarshal(src, &null)
		if err == nil {
			err = errors.New("invalid JSON")
		}
		return err
	}
	r := bytes.NewReader(src)
	var node goyaml.Node
	if err := goyaml.NewDecoder(r).Decode(&node); err != nil {
		return err
	}
	dec := &Decoder{}
	out, err := dec.Decode(&node)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(v).Elem()
	if !val.CanSet() {
		return errors.New("cannot set value")
	}
	val.Set(reflect.ValueOf(out))
	return nil
}
