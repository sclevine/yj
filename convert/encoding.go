package convert

import (
	"io"
	"math"
)

type Encoding interface {
	String() string
	Encode(w io.Writer, in interface{}) error
	Decode(r io.Reader) (interface{}, error)
}

type SpecialFloats int

const (
	FloatsReal SpecialFloats = iota
	FloatsString
	FloatsNumber
)

func (s SpecialFloats) NaN() interface{} {
	switch s {
	case FloatsReal:
		return math.NaN()
	case FloatsString:
		return "NaN"
	case FloatsNumber:
		return (*float64)(nil)
	}
	panic("NaN: invalid special float type")
}

func (s SpecialFloats) PosInf() interface{} {
	switch s {
	case FloatsReal:
		return math.Inf(1)
	case FloatsString:
		return "Infinity"
	case FloatsNumber:
		return math.MaxFloat64
	}
	panic("PosInf: invalid special float type")
}

func (s SpecialFloats) NegInf() interface{} {
	switch s {
	case FloatsReal:
		return math.Inf(-1)
	case FloatsString:
		return "-Infinity"
	case FloatsNumber:
		return -math.MaxFloat64
	}
	panic("NegInf: invalid special float type")
}
