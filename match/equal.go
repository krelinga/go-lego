package match

import "fmt"

func Equal[T comparable](expected T) Matcher {
	return NewFunc(MetaHere(), func(result *Result, actual T) {
		if actual == expected {
			result.Accepted = true
		}
		result.Why = func() string {
			var outcome string
			if result.Accepted {
				outcome = "=="
			} else {
				outcome = "!="
			}
			return fmt.Sprintf("%#v %s %#v", actual, outcome, expected)
		}()
	})
}
