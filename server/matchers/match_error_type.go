package matchers

import (
	"github.com/onsi/gomega/types"

	"fmt"
	"reflect"
)

func MatchErrorType(expected interface{}) types.GomegaMatcher {
	return &matchErrorTypeMatcher{
		expected: reflect.TypeOf(expected),
	}
}

type matchErrorTypeMatcher struct {
	expected reflect.Type
}

func (matcher *matchErrorTypeMatcher) Match(actual interface{}) (success bool, err error) {
	return (reflect.TypeOf(actual) == matcher.expected), nil
}

func (matcher *matchErrorTypeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto be of type\n\t%#v", actual, matcher.expected)
}

func (matcher *matchErrorTypeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be of type\n\t%#v", actual, matcher.expected)
}
