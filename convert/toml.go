package convert

import (
	"bytes"
	"crypto/rand"
	"io"
	"strconv"

	gotoml "github.com/BurntSushi/toml"

	"github.com/sclevine/yj/v5/order"
	"github.com/sclevine/yj/v5/toml"
)

type TOML struct {
	SpecialFloats
	Indent bool
}

// Reuse the prefix so repeated conversions reuse reflect.StructOf types.
var tomlKeyPrefix = "_yj_" + rand.Text() + "_"

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) error {
	var buf bytes.Buffer
	tomlEnc := gotoml.NewEncoder(&buf)
	if !t.Indent {
		tomlEnc.Indent = ""
	}
	enc := toml.Encoder{
		NaN:    t.NaN(),
		PosInf: t.PosInf(),
		NegInf: t.NegInf(),
	}
	// Struct tags interpret commas, empty names and "-" instead of preserving keys.
	keys := make(map[string]string)
	in = aliasTOMLKeys(in, tomlKeyPrefix, keys)
	out, err := enc.Encode(in)
	if err != nil {
		return err
	}
	if err := tomlEnc.Encode(out); err != nil {
		return err
	}
	return restoreTOMLKeys(w, buf.Bytes(), keys)
}

func aliasTOMLKeys(in interface{}, prefix string, keys map[string]string) interface{} {
	switch in := in.(type) {
	case order.MapSlice:
		out := make(order.MapSlice, len(in))
		for i, item := range in {
			if key, ok := item.Key.(string); ok {
				alias := prefix + strconv.Itoa(len(keys))
				keys[alias] = gotoml.Key{key}.String()
				item.Key = alias
			}
			item.Val = aliasTOMLKeys(item.Val, prefix, keys)
			out[i] = item
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(in))
		for i, item := range in {
			out[i] = aliasTOMLKeys(item, prefix, keys)
		}
		return out
	default:
		return in
	}
}

func restoreTOMLKeys(w io.Writer, data []byte, keys map[string]string) error {
	var out bytes.Buffer
	for i := 0; i < len(data); {
		start := i
		switch data[i] {
		case '"', '\'':
			quote := data[i : i+1]
			if bytes.HasPrefix(data[i:], bytes.Repeat(quote, 3)) {
				quote = data[i : i+3]
			}
			i += len(quote)
			for i < len(data) {
				if data[i] == '\\' && quote[0] == '"' {
					i += 2
					if i > len(data) {
						i = len(data)
					}
				} else if bytes.HasPrefix(data[i:], quote) {
					i += len(quote)
					if len(quote) == 3 {
						for n := 0; n < 2 && i < len(data) && data[i] == quote[0]; n++ {
							i++
						}
					}
					break
				} else {
					i++
				}
			}
		case '#':
			if end := bytes.IndexByte(data[i:], '\n'); end >= 0 {
				i += end + 1
			} else {
				i = len(data)
			}
		default:
			if !tomlBareKeyByte(data[i]) {
				i++
				break
			}
			for i < len(data) && tomlBareKeyByte(data[i]) {
				i++
			}
			if key, ok := keys[string(data[start:i])]; ok {
				out.WriteString(key)
				continue
			}
		}
		out.Write(data[start:i])
	}
	n, err := w.Write(out.Bytes())
	if err == nil && n != out.Len() {
		return io.ErrShortWrite
	}
	return err
}

func tomlBareKeyByte(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '_' || b == '-'
}

func (t TOML) Decode(r io.Reader) (interface{}, error) {
	var out interface{}
	md, err := gotoml.NewDecoder(r).Decode(&out)
	if err != nil {
		return nil, err
	}
	dec := toml.Decoder{
		NaN:    t.NaN(),
		PosInf: t.PosInf(),
		NegInf: t.NegInf(),
	}
	return dec.Decode(out, md.Keys())
}
