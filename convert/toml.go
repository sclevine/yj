package convert

import (
	"io"

	gotoml "github.com/BurntSushi/toml"

	"github.com/sclevine/yj/v5/toml"
)

type TOML struct {
	SpecialFloats
	Indent bool
}

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) error {
	tomlEnc := gotoml.NewEncoder(w)
	if !t.Indent {
		tomlEnc.Indent = ""
	}
	enc := toml.Encoder{
		NaN:    t.NaN(),
		PosInf: t.PosInf(),
		NegInf: t.NegInf(),
	}
	out, err := enc.Encode(in)
	if err != nil {
		return err
	}
	return tomlEnc.Encode(out)
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
