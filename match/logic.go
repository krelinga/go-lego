package match

import (
	"fmt"
	"reflect"
)

func Not(m Matcher) Matcher {
	return NewValFunc(MetaHere(), func(result *Result, val reflect.Value) {
		childRes := m.Match(val)
		result.Children = append(result.Children, Child{
			Result: childRes,
		})
		if childRes.Err != nil {
			result.Err = fmt.Errorf("child matcher failed")
			return
		}
		result.Accepted = !childRes.Accepted
		result.Explain = func() string {
			outcomeStr := func(accepted bool) string {
				if accepted {
					return "accepted"
				}
				return "rejected"
			}
			childOutcome := outcomeStr(childRes.Accepted)
			notOutcome := outcomeStr(result.Accepted)
			return fmt.Sprintf("child result was %s, so NOT result is %s", childOutcome, notOutcome)
		}
	})
}

func AllOf(matchers ...Matcher) Matcher {
	return NewValFunc(MetaHere(), func(result *Result, val reflect.Value) {
		result.Children = make([]Child, len(matchers))
		accept := true
		for i, m := range matchers {
			childRes := m.Match(val)
			result.Children[i] = Child{
				Name:  fmt.Sprintf("child %d", i),
				Result: childRes,
			}
			if childRes.Err != nil && result.Err == nil {
				result.Err = fmt.Errorf("child matcher %d failed", i)
			}
			if !childRes.Accepted {
				accept = false
			}
		}
		if result.Err == nil && accept {
			result.Accepted = true
		}
		result.Explain = func() string {
			if result.Accepted {
				return "all child matchers accepted"
			}
			return "some child matchers rejected"
		}
	})
}