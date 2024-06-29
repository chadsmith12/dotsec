package passbolt_test

import (
	"reflect"
	"testing"

	"github.com/chadsmith12/dotsec/secrets"
)

type testCase struct {
	name     string
	input    []string
	expected []secrets.SecretData
}

func TestParsingSecrets(t *testing.T) {
	testCases := createParseTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := secrets.SecretDataFromSlice(testCase.input)
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
			expected: []secrets.SecretData{{Key: "ExampleKey", Value: "MyData"}},
		},
		{
			name:     "Returns Empty on Empty Slice",
			input:    []string{},
			expected: []secrets.SecretData{},
		},
		{
			name:     "Ignores invalid separators",
			input:    []string{"BadKey:BadData"},
			expected: []secrets.SecretData{},
		},
	}
}
