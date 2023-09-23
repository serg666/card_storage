package internal

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type validateMethods struct {
}

func (vm *validateMethods) json(dst interface{}, src io.Reader) httpError {
	if src == nil {
		return httpErrValidationFailed
	}

	if err := json.NewDecoder(src).Decode(dst); err != nil {
		logrus.Warnf("Validate by json failed: src=%v, err=%v", src, err)
		return httpErrValidationFailed
	}

	return nil
}

func (vm *validateMethods) structs(s interface{}) httpError {
	newValidator := validator.New()
	err := newValidator.RegisterValidation("notempty", func(fl validator.FieldLevel) bool {
		return fl.Field().Len() != 0
	})
	if err != nil {
		logrus.Warnf("Register validation failed, err=%v", err)
		return httpErrValidationFailed
	}

	if err := newValidator.Struct(s); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return httpValidationErrorByValidationErrors(errs)
		}
		logrus.Warnf("Validate struct failed: struct=%v, err=%v", s, err)
		return httpErrValidationFailed
	}

	return nil
}
