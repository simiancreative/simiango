package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/validator.v2"
)

type resource struct {
	Wibble string `json:"wibble" validate:"nonzero"`
}

func TestValidateSuccess(t *testing.T) {
	r := resource{Wibble: "23"}
	result := Validate(r)

	assert.Equal(t, false, result.HasErrors())
}

func TestValidateFailed(t *testing.T) {
	r := resource{}
	result := Validate(r)

	assert.Equal(t, true, result.HasErrors())
	assert.Equal(t, "Wibble: zero value", result.Errors().Error())
}

func TestValidateResultError(t *testing.T) {
	r := resource{}
	result := Validate(r)

	err := result.ResultError()
	assert.Equal(t, "request_validation_failed", err.GetMessage())
}

func TestValidateError(t *testing.T) {
	r := resource{}
	result := Validate(r)

	err := result.Errors()
	errs := err.(validator.ErrorMap)

	assert.Equal(t, "zero value", errs["Wibble"][0].(validator.TextErr).Error())
}
