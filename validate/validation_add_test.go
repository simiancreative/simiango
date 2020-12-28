package validate

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Very simple validation func
func notZZ(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("unsupported_type")
	}
	if st.String() == "ZZ" {
		return errors.New("value cannot be ZZ")
	}
	return nil
}

type sample struct {
	A string `validate:"nonzero,notzz"`
}

func TestAddValidation(t *testing.T) {
	AddValidation("notzz", notZZ)

	s := sample{"ZZ"}
	result := Validate(s)

	assert.Equal(t, "A: value cannot be ZZ", result.Errors().Error())
}
