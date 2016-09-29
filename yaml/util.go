package yaml

import "fmt"

func catchFailure(err *error) {
	if r := recover(); r != nil {
		var ok bool
		if *err, ok = r.(error); !ok {
			*err = fmt.Errorf("unexpected failure: %v", r)
		}
	}
}
