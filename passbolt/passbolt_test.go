package passbolt_test

import (
	"reflect"
	"testing"

	"github.com/chadsmith12/dotsec/passbolt"
)

type testCase struct {
	name     string
	input    []string
	expected []passbolt.SecretData
}

func TestParsingSecrets(t *testing.T) {
	testCases := createParseTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := passbolt.SecretDataFromSlice(testCase.input)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("SecretDataFromSlice() returned %v, expected %v\n", actual, testCase.expected)
			}
		})
	}
}

func createParseTestCases() []testCase {
	return []testCase{
		{
			name:     "Parses Single Line",
			input:    []string{"ExampleKey = MyData"},
			expected: []passbolt.SecretData{{Key: "ExampleKey", Value: "MyData"}},
		},
		{
			name:     "Returns Empty on Empty Slice",
			input:    []string{},
			expected: []passbolt.SecretData{},
		},
		{
			name:     "Ignores invalid separators",
			input:    []string{"BadKey:BadData"},
			expected: []passbolt.SecretData{},
		},
	}
}
