package nulls

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Marshaller interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

type test struct {
	Name            string
	GetInst         func() interface{}
	Valid           func(interface{}) bool
	Value           func(interface{}) interface{}
	Matcher         interface{}
	Param           interface{}
	MarshalledParam interface{}
	FailParam       interface{}
}

var tests = []test{}

func registerTest(t test) {
	tests = append(tests, t)
}

func TestScan(t *testing.T) {
	for _, test := range tests {
		inst := test.GetInst()
		err := inst.(sql.Scanner).Scan(test.Param)

		assert.NoError(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(t, test.Valid(inst), true, fmt.Sprintf("%v should be equal", test.Name))
		assert.Equal(t, test.Value(inst), test.Matcher, fmt.Sprintf("%v should be equal", test.Name))
	}
}

func TestScanFail(t *testing.T) {
	for _, test := range tests {
		inst := test.GetInst()
		err := inst.(sql.Scanner).Scan(test.FailParam)

		assert.Error(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(
			t,
			err.Error(),
			fmt.Sprintf("value not valid: %T, %v", test.FailParam, test.FailParam),
			fmt.Sprintf("%v should have error", test.Name),
		)
		assert.Equal(
			t,
			test.Valid(inst),
			false,
			fmt.Sprintf("%v shouldnt be valid", test.Name),
		)
	}
}

func TestMarshall(t *testing.T) {
	for _, test := range tests {
		inst := test.GetInst()
		err := inst.(sql.Scanner).Scan(test.Param)
		bytes, _ := inst.(Marshaller).MarshalJSON()
		expectedBytes, _ := json.Marshal(test.Value(inst))

		assert.NoError(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(t, bytes, expectedBytes, fmt.Sprintf("%v should be equal", test.Name))
	}
}

func TestUnMarshall(t *testing.T) {
	for _, test := range tests {
		inst := test.GetInst()
		expectedBytes, _ := json.Marshal(test.MarshalledParam)
		err := inst.(Marshaller).UnmarshalJSON(expectedBytes)

		assert.NoError(t, err, fmt.Sprintf("%v shouldnt have error", test.Name))
		assert.Equal(t, test.Value(inst), test.Matcher, fmt.Sprintf("%v should be equal", test.Name))
	}
}
