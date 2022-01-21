package nulls

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	Name      string
	Func      func(interface{}) (bool, interface{}, error)
	Param     interface{}
	Matcher   interface{}
	FailParam interface{}
}

var tests = []test{}

func registerTest(t test) {
	tests = append(tests, t)
}

func TestScan(t *testing.T) {
	for _, test := range tests {
		valid, value, err := test.Func(test.Param)

		assert.NoError(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(t, valid, true, fmt.Sprintf("%v should be equal", test.Name))
		assert.Equal(t, value, test.Matcher, fmt.Sprintf("%v should be equal", test.Name))
	}
}

func TestScanFail(t *testing.T) {
	for _, test := range tests {
		valid, _, err := test.Func(test.FailParam)

		assert.Error(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(
			t,
			err.Error(),
			fmt.Sprintf("value not valid: %T, %v", test.FailParam, test.FailParam),
			fmt.Sprintf("%v should have error", test.Name),
		)
		assert.Equal(
			t,
			valid,
			false,
			fmt.Sprintf("%v shouldnt be valid", test.Name),
		)
	}
}
