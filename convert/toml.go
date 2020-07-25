package convert

import (
	"io"

	gotoml "github.com/pelletier/go-toml"

	"github.com/sclevine/yj/toml"
)

type TOML struct {
	SpecialFloats
	Indent bool
}

func (TOML) String() string {
	return "TOML"
}

func (t TOML) Encode(w io.Writer, in interface{}) error {
	tomlEnc := gotoml.NewEncoder(&trimWriter{w: w})
	tomlEnc.Order(gotoml.OrderPreserve)
	if !t.Indent {
		tomlEnc.Indentation("")
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

type trimWriter struct {
	w    io.Writer
	done bool
}

func (w *trimWriter) Write(p []byte) (n int, err error) {
	trimmed := false
	if !w.done && len(p) > 0 && p[0] == '\n' {
		p = p[1:]
		trimmed = true
	}
	n, err = w.w.Write(p)
	if (trimmed && err == nil) || n > 0 {
		w.done = true
		if trimmed {
			n++
		}
	}
	return n, err
}

func (t TOML) Decode(r io.Reader) (interface{}, error) {
	tree, err := gotoml.LoadReader(r)
	if err != nil {
		return nil, err
	}
	dec := toml.Decoder{
		NaN:    t.NaN(),
		PosInf: t.PosInf(),
		NegInf: t.NegInf(),
	}
	return dec.Decode(tree)
}
