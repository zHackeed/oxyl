package validator

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	validator *validator.Validate
}

func NewStructValidator(customValidators ...*validator.Validate) *StructValidator {
	var parser *validator.Validate

	if len(customValidators) > 0 {
		parser = customValidators[0]
	} else {
		parser = validator.New()
	}

	return &StructValidator{
		validator: parser,
	}
}

func (v *StructValidator) Validate(data any) error {
	return v.validator.Struct(data)
}
