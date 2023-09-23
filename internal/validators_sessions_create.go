package internal

import (
	"io"

	"github.com/serg666/card_storage/models"
)

type validatorSessionsCreate struct {
	vm   *validateMethods
	Key  string `json:"key" validate:"required"`
	Body string `json:"body" validate:"required"`
}

func (v *validatorSessionsCreate) validate(src io.Reader) httpError {
	if err := v.vm.json(v, src); err != nil {
		return err
	}

	return v.vm.structs(v)
}

func (v *validatorSessionsCreate) toSessionModel(ckey *models.Key) *models.Session {
	return models.NewSession(v.Key, v.Body, ckey)
}

func (vm *validateMethods) SessionsCreate(src io.Reader) (*validatorSessionsCreate, httpError) {
	v := &validatorSessionsCreate{
		vm: vm,
	}
	return v, v.validate(src)
}
