package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	parser "github.com/v1adhope/query-parser"
)

func TestParseExpressiontToJSON(t *testing.T) {
	testCases := []struct {
		Name        string
		Input       string
		Expected    string
		ExpectedErr error
	}{
		{
			Name:        "Empty input",
			Input:       "",
			Expected:    "",
			ExpectedErr: parser.ERR_EMPTY_RAW,
		},
		{
			Name:     "There're not logic operations",
			Input:    "drunk=coffee&expression=milk==true",
			Expected: `{"fields":{"drunk":"coffee"}, "expression":{"logicOperation":{"Single":[{"key":{"milk":[{"value":"true", "operation":"=="}]}}]}}}`,
		},
		{
			Name:     "There's not expression",
			Input:    "name=coffee&milk=70%",
			Expected: `{"fields":{"milk":"70%", "name":"coffee"}}`,
		},
		{
			Name:     "There's only expression",
			Input:    "expression=ip.src == 1.1.1.1 || ip.src == 62.13.171.36 && port.src != 10000",
			Expected: `{"expression":{"logicOperation":{"\u0026\u0026":[{"key":{"port.src":[{"value":"10000", "operation":"!="}]}}], "||":[{"key":{"ip.src":[{"value":"1.1.1.1", "operation":"=="}]}}, {"key":{"ip.src":[{"value":"62.13.171.36", "operation":"=="}]}}]}}}`,
		},
	}

	for _, tc := range testCases {
		SUT, err := parser.ParseRawToJSON(tc.Input)

		assert.Equal(t, tc.Expected, SUT, tc.Name)
		assert.Equal(t, tc.ExpectedErr, err, tc.Name)
	}
}
