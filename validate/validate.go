package validate

import (
	"github.com/simiancreative/simiango/service"

	"gopkg.in/validator.v2"
)

func Valid(v interface{}, tags string) error {
	return validator.Valid(v, tags)
}

func Validate(v interface{}) Validator {
	errs := validator.Validate(v)
	return Validator{Errs: errs}
}

func AddValidation(name string, f validator.ValidationFunc) error {
	return validator.SetValidationFunc(name, f)
}

type Validator struct {
	Errs error
}

func (dv *Validator) HasErrors() bool {
	return dv.Errs != nil
}

func (dv *Validator) Errors() error {
	return dv.Errs
}

func (dv *Validator) ResultError() *service.ResultError {
	if dv.Errs == nil {
		return nil
	}

	reasons := []map[string]interface{}{}

	for k, v := range dv.Errs.(validator.ErrorMap) {
		reasons = append(reasons, map[string]interface{}{})
		reasons[0][k] = v
	}

	message := "request_validation_failed"
	return &service.ResultError{
		Status:     422,
		ErrMessage: message,
		Message:    message,
		Reasons:    reasons,
	}
}
