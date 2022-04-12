package toml

import (
	"fmt"
	"math"

	gotoml "github.com/BurntSushi/toml"

	"github.com/sclevine/yj/v5/order"
)

type Decoder struct {
	// If set, NaN, Inf, etc. are replaced by the set values
	NaN, PosInf, NegInf interface{}
}

func (d *Decoder) Decode(toml interface{}, keys []gotoml.Key) (normal interface{}, err error) {
	defer catchFailure(&err)
	v, _ := d.decode(toml, keys, nil, 0)
	return v, nil
}

func (d Decoder) decode(v interface{}, keys []gotoml.Key, key gotoml.Key, pos int) (interface{}, int) {
	switch v := v.(type) {
	case map[string]interface{}:
		ks, npos := uniqueKeys(keys, key, pos, len(v))
		if len(ks) != len(v) {
			panic(fmt.Errorf("key mismatch, %d vs. %d", len(ks), len(v)))
		}
		out := make(order.MapSlice, 0, len(ks))
		for _, k := range ks {
			next, ok := v[k]
			if !ok {
				panic(fmt.Errorf("missing key `%s'", k))
			}
			val, dpos := d.decode(next, keys, append(key, k), pos)
			if dpos > npos {
				npos = dpos
			}
			out = append(out, order.MapItem{Key: k, Val: val})
		}
		if len(keys) > npos && keysEqual(keys[npos], key) {
			npos++
		}
		return out, npos
	case []map[string]interface{}:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			var val interface{}
			val, pos = d.decode(item, keys, key, pos)
			out = append(out, val)
		}
		return out, pos
	case []interface{}:
		out := make([]interface{}, 0, len(v))
		for _, item := range v {
			var val interface{}
			val, pos = d.decode(item, keys, key, pos)
			out = append(out, val)
		}
		return out, pos
	default:
		return d.convert(v), pos
	}
}

func uniqueKeys(keys []gotoml.Key, prefix gotoml.Key, pos, n int) ([]string, int) {
	m := make(map[string]struct{})
	var out []string
	end := pos
	for i, k := range keys[pos:] {
		if n == 0 {
			break
		}
		rest, ok := startsWith(k, prefix)
		if !ok {
			continue
		}
		if len(rest) == 0 {
			end = pos + i + 1 // actually needed?
			continue
		}
		r := rest[0]
		if _, ok := m[r]; !ok {
			m[r] = struct{}{}
			out = append(out, r)
			n--
			end = pos + i + 1
		}
	}
	return out, end
}

func (d Decoder) convert(in interface{}) interface{} {
	switch in := in.(type) {
	case float64:
		switch {
		case d.NaN != nil && math.IsNaN(in):
			return d.NaN
		case d.PosInf != nil && math.IsInf(in, 1):
			return d.PosInf
		case d.NegInf != nil && math.IsInf(in, -1):
			return d.NegInf
		}
		return in
	}
	return in
}

func startsWith(key, prefix gotoml.Key) (rest []string, ok bool) {
	if len(key) < len(prefix) {
		return nil, false
	}
	for i := range prefix {
		if key[i] != prefix[i] {
			return nil, false
		}
	}
	return key[len(prefix):], true
}

func keysEqual(k1, k2 gotoml.Key) bool {
	if len(k1) != len(k2) {
		return false
	}
	for i := range k1 {
		if k1[i] != k2[i] {
			return false
		}
	}
	return true
}
