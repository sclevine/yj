package convert

import (
	"bytes"
	"io"
	"io/ioutil"

	gohcl "github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	hcljson "github.com/hashicorp/hcl/json/parser"

	"github.com/sclevine/yj/hcl"
)

type HCL struct{}

func (HCL) String() string {
	return "HCL"
}

func (HCL) Encode(w io.Writer, in interface{}) error {
	j := &bytes.Buffer{}
	if err := (JSON{}).Encode(j, in); err != nil {
		return err
	}
	f, err := hcljson.Parse(j.Bytes())
	if err != nil {
		return err
	}
	return printer.Fprint(w, f)
}

func (HCL) Decode(r io.Reader) (interface{}, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f, err := gohcl.ParseBytes(in)
	if err != nil {
		return nil, err
	}
	dec := hcl.Decoder{}
	return dec.Decode(f.Node)
}
