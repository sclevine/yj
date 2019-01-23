package convert

import "io"

type Encoding interface {
	String() string
	Encode(w io.Writer, in interface{}) error
	Decode(r io.Reader) (interface{}, error)
}
