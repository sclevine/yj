package toml

import (
	"math"

	"github.com/sclevine/yj/v5/order"
)

type Encoder struct {
	// If set, the set values will be converted to NaN, Inf, etc.
	NaN, PosInf, NegInf interface{}
}

func (e *Encoder) Encode(normal interface{}) (toml interface{}, err error) {
	defer catchFailure(&err)
	return e.encode(normal), nil
}

func (e *Encoder) encode(val interface{}) interface{} {
	switch val := val.(type) {
	case order.MapSlice:
		for i, item := range val {
			val[i].Val = e.encode(item.Val)
		}
		s, err := val.Struct()
		if err != nil {
			panic(err)
		}
		return s
	case []interface{}:
		out := make([]interface{}, 0, len(val))
		for _, v := range val {
			if v == nil {
				continue
			}
			out = append(out, e.encode(v))
		}
		return out
	default:
		return e.convert(val)
	}
}

func (e *Encoder) convert(in interface{}) interface{} {
	switch in {
	case nil:
		return nil
	case e.NaN:
		return math.NaN()
	case e.PosInf:
		return math.Inf(1)
	case e.NegInf:
		return math.Inf(-1)
	}
	return in
}
