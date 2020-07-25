package yaml

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	goyaml "gopkg.in/yaml.v3"

	"github.com/sclevine/yj/order"
)

// NOTE: some parts copied from gopkg.in/yaml.v3

const (
	// 400,000 decode operations is ~500kb of dense object declarations, or
	// ~5kb of dense object declarations with 10000% alias expansion
	aliasRatioRangeLow = 400000

	// 4,000,000 decode operations is ~5MB of dense object declarations, or
	// ~4.5MB of dense object declarations with 10% alias expansion
	aliasRatioRangeHigh = 4000000

	// aliasRatioRange is the range over which we scale allowed alias ratios
	aliasRatioRange = float64(aliasRatioRangeHigh - aliasRatioRangeLow)

	longTagPrefix = "tag:yaml.org,2002:"
	mergeTag      = "!!merge"
)

var ErrNotMaps = errors.New("map merge requires map or sequence of maps as the value")

type Decoder struct {
	KeyMarshal func(interface{}) ([]byte, error)

	// If set, NaN, Inf, etc. are replaced by the set values
	NaN, PosInf, NegInf          interface{}
	KeyNaN, KeyPosInf, KeyNegInf interface{}
}

// Decode decodes a YAML node tree into an the normalized object format.
func (d *Decoder) Decode(node *goyaml.Node) (normal interface{}, err error) {
	defer catchFailure(&err)
	dt := decodeTracker{Decoder: d}
	return dt.normalize(node), nil
}

type decodeTracker struct {
	*Decoder
	aliases     map[*goyaml.Node]struct{}
	doc         *goyaml.Node
	decodeCount int
	aliasCount  int
	aliasDepth  int
}

func (d *decodeTracker) normalize(n *goyaml.Node) interface{} {
	d.decodeCount++
	if d.aliasDepth > 0 {
		d.aliasCount++
	}
	if d.aliasCount > 100 && d.decodeCount > 1000 && float64(d.aliasCount)/float64(d.decodeCount) > allowedAliasRatio(d.decodeCount) {
		panic(fmt.Errorf("document contains excessive aliasing"))
	}

	switch n.Kind {
	case goyaml.DocumentNode:
		return d.document(n)
	case goyaml.AliasNode:
		return d.alias(n)
	case goyaml.ScalarNode:
		var out interface{}
		if err := n.Decode(&out); err != nil {
			panic(fmt.Errorf("scalar decode error: %s", err))
		}
		switch out := out.(type) {
		case float64:
			return d.float(out)
		}
		return d.other(out)
	case goyaml.MappingNode:
		return d.mapping(n)
	case goyaml.SequenceNode:
		return d.sequence(n)
	case 0:
		if n.IsZero() {
			return nil
		}
		fallthrough
	default:
		panic(fmt.Errorf("cannot decode node with unknown kind %d", n.Kind))
	}
}

func (d *decodeTracker) other(in interface{}) interface{} {
	switch reflect.ValueOf(in).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Float32:
		panic(fmt.Errorf("unexpected type: %#v", in))
	}
	return in
}

func (d *decodeTracker) float(in float64) interface{} {
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

func (d *decodeTracker) key(n *goyaml.Node) string {
	// Decoder remains reentrant, but decodeTracker need not be
	defer func(dec *Decoder) {
		d.Decoder = dec
	}(d.Decoder)
	kdec := *d.Decoder
	kdec.NaN = d.KeyNaN
	kdec.PosInf = d.KeyPosInf
	kdec.NegInf = d.KeyNegInf
	d.Decoder = &kdec
	switch key := d.normalize(n).(type) {
	case string:
		return key
	case fmt.Stringer:
		return key.String()
	default:
		out, err := d.KeyMarshal(key)
		if err != nil {
			panic(err)
		}
		return string(out)
	}
}

func allowedAliasRatio(decodeCount int) float64 {
	switch {
	case decodeCount <= aliasRatioRangeLow:
		// allow 99% to come from alias expansion for small-to-medium documents
		return 0.99
	case decodeCount >= aliasRatioRangeHigh:
		// allow 10% to come from alias expansion for very large documents
		return 0.10
	default:
		// scale smoothly from 99% down to 10% over the range.
		// this maps to 396,000 - 400,000 allowed alias-driven decodes over the range.
		// 400,000 decode operations is ~100MB of allocations in worst-case scenarios (single-item maps).
		return 0.99 - 0.89*(float64(decodeCount-aliasRatioRangeLow)/aliasRatioRange)
	}
}

func (d *decodeTracker) document(n *goyaml.Node) interface{} {
	if len(n.Content) != 1 {
		panic(fmt.Errorf("invalid document"))
	}
	d.doc = n
	return d.normalize(n.Content[0])
}

func (d *decodeTracker) alias(n *goyaml.Node) interface{} {
	if d.aliases == nil {
		d.aliases = make(map[*goyaml.Node]struct{})
	}
	if _, ok := d.aliases[n]; ok {
		// TODO this could actually be allowed in some circumstances.
		panic(fmt.Errorf("anchor '%s' value contains itself", n.Value))
	}
	d.aliases[n] = struct{}{}
	d.aliasDepth++
	out := d.normalize(n.Alias)
	d.aliasDepth--
	delete(d.aliases, n)
	return out
}

func (d *decodeTracker) sequence(n *goyaml.Node) []interface{} {
	out := make([]interface{}, 0, len(n.Content))
	for _, c := range n.Content {
		out = append(out, d.normalize(c))
	}
	return out
}

func shortTag(tag string) string {
	// TODO This can easily be made faster and produce less garbage.
	if strings.HasPrefix(tag, longTagPrefix) {
		return "!!" + tag[len(longTagPrefix):]
	}
	return tag
}

func isMerge(n *goyaml.Node) bool {
	return n.Kind == goyaml.ScalarNode && n.Value == "<<" && (n.Tag == "" || n.Tag == "!" || shortTag(n.Tag) == mergeTag)
}

func (d *decodeTracker) mapping(n *goyaml.Node) order.MapSlice {
	l := len(n.Content)
	out := make(order.MapSlice, 0, l/2)

	for i := 0; i < l; i += 2 {
		if isMerge(n.Content[i]) {
			out = d.merge(out, n.Content[i+1])
			continue
		}
		out = append(out, order.MapItem{
			Key: d.key(n.Content[i]),
			Val: d.normalize(n.Content[i+1]),
		})
	}
	return out
}

func (d *decodeTracker) merge(m order.MapSlice, n *goyaml.Node) order.MapSlice {
	switch n.Kind {
	case goyaml.AliasNode:
		if n.Alias != nil && n.Alias.Kind != goyaml.MappingNode {
			panic(ErrNotMaps)
		}
		fallthrough
	case goyaml.MappingNode:
		in, ok := d.normalize(n).(order.MapSlice)
		if !ok {
			panic(ErrNotMaps)
		}
		return m.Merge(in)
	case goyaml.SequenceNode:
		// Step backwards as earlier nodes take precedence.
		for i := len(n.Content) - 1; i >= 0; i-- {
			ni := n.Content[i]
			if ni.Kind == goyaml.AliasNode {
				if ni.Alias != nil && ni.Alias.Kind != goyaml.MappingNode {
					panic(ErrNotMaps)
				}
			} else if ni.Kind != goyaml.MappingNode {
				panic(ErrNotMaps)
			}
			in, ok := d.normalize(n).(order.MapSlice)
			if !ok {
				panic(ErrNotMaps)
			}
			m = m.Merge(in)
		}
		return m
	default:
		panic(ErrNotMaps)
	}
}
