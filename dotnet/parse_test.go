package dotnet_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/chadsmith12/dotsec/dotnet"
)

type testCase struct {
	name     string
	input    string
	expected []string
}

func TestSecretParsing(t *testing.T) {
	testCases := createParseTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buf := bytes.NewBufferString(testCase.input)
			result, err := dotnet.ParseSecrets(*buf)
			if err != nil {
				t.Errorf("ParseSecrets() error = %v, expected %v\n", err, testCase.expected)
			}

			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("ParseSecrets() result = %v, expected %v\n", result, testCase.expected)
			}
		})
	}
}

func createParseTestCases() []testCase {
	return []testCase{
		{
			name:     "Parses Valid Input",
			input:    "ExampleKey = MyData",
			expected: []string{"ExampleKey = MyData"},
		},
		{
			name:     "Parses No Secrets",
			input:    "No secrets configured",
			expected: []string{},
		},
	}
}
