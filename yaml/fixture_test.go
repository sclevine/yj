package yaml_test

import "math"

type stringer string

func (s stringer) String() string {
	return "stringer=" + string(s)
}

type F struct {
	F string `json:"f"`
}

var yamlFixture = []interface{}{
	map[interface{}]interface{}{
		1.1:           []interface{}{"a"},
		"b":           []interface{}{2.2},
		stringer("c"): map[interface{}]interface{}{3.3: 4.4},
		1:             F{"not NaN"},
		math.Inf(1):   math.Inf(-1),
		math.Inf(-1):  math.Inf(1),
	},
}

var jsonFixture = []interface{}{
	map[string]interface{}{
		"1.1":               []interface{}{"a"},
		"b":                 []interface{}{2.2},
		"stringer=c":        map[string]interface{}{"3.3": 4.4},
		"1":                 F{"not NaN"},
		`{"f":"Infinity"}`:  F{"-Infinity"},
		`{"f":"-Infinity"}`: F{"Infinity"},
	},
}
