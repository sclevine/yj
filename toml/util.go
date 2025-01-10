package toml

import "fmt"

// useful for catching panics in external packages
func catchFailure(err *error) {
	if r := recover(); r != nil {
		var ok bool
		if *err, ok = r.(error); !ok {
			*err = fmt.Errorf("unexpected failure: %v", r)
		}
	}
}
