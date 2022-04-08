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
	//fmt.Printf("KEYS %+v\n", keys)
	v, _ := d.decode(toml, keys, nil, 0)
	return v, nil
}

func (d Decoder) decode(v interface{}, keys []gotoml.Key, key gotoml.Key, pos int) (interface{}, int) {
	switch v := v.(type) {
	case map[string]interface{}:
		ks, npos := uniqueKeys(keys, key, pos, len(v))
		//fmt.Printf("DECODE MAP [key:] %+v [keys:] %+v [l] %d [p] %d [np] %d\n", key, ks, len(v), pos, npos)
		//fmt.Printf("MAP THING %+v\n", v)
		if len(ks) != len(v) {
			panic(fmt.Errorf("key mismatch, %d vs. %d", len(ks), len(v)))
		}
		var m order.MapSlice
		for _, k := range ks {
			next, ok := v[k]
			if !ok {
				panic(fmt.Errorf("missing key `%s'", k))
			}
			val, _ := d.decode(next, keys, append(key, k), pos)
			m = append(m, order.MapItem{Key: k, Val: val})
		}
		return m, npos
	case []map[string]interface{}:
		//fmt.Printf("DECODE LIST %+v %d %d\n", key, pos, len(v))
		var out []interface{}
		for _, item := range v {
			var val interface{}
			val, pos = d.decode(item, keys, key, pos)
			out = append(out, val)
		}
		return out, 0
	case []interface{}:
		//fmt.Printf("DECODE LIST %+v %d %d\n", key, pos, len(v))
		var out []interface{}
		for _, item := range v {
			var val interface{}
			val, pos = d.decode(item, keys, key, pos)
			out = append(out, val)
		}
		return out, 0
	default:
		return d.convert(v), 0
	}
}

func uniqueKeys(keys []gotoml.Key, key gotoml.Key, pos, n int) ([]string, int) {
	m := make(map[string]struct{})
	var out []string
	end := 0
	seen := false
	for i, k := range keys[pos:] {
		rest, ok := startsWith(k, key)
		if !ok {
			continue
		}
		if len(rest) == 0 {
			if seen {
				end = i + pos // needed?
				break
			}
			seen = true
			continue
		}
		r := rest[0]
		if seen {
			if _, ok := m[r]; !ok {
				//if n == 0 {
				//	break
				//}
				m[r] = struct{}{}
				out = append(out, r)
				//n--
			}
			end = i + pos + 1
		} else {
			if n == 0 {
				break
			}
			if _, ok := m[r]; !ok {
				m[r] = struct{}{}
				out = append(out, r)
				n--
				end = i + pos + 1
			}
		}
	}
	return out, end
}

func uniqueKeysOld(keys []gotoml.Key, key gotoml.Key, pos, n int) ([]string, int) {
	// seeing name of table switches to greedy alg, otherwise stop at last new key?
	m := make(map[string]struct{})
	var out []string
	end := 0
	for i, k := range keys[pos:] {
		rest, ok := startsWith(k, key)
		if !ok || len(rest) == 0 {
			continue
		}
		r := rest[0]
		if _, ok := m[r]; !ok {
			if n == 0 {
				break
			}
			m[r] = struct{}{}
			out = append(out, r)
			n--
		}
		end = i + pos + 1
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
