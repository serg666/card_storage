package internal

import (
	"github.com/gin-gonic/gin"
)

type validatorCardByToken struct {
	vm    *validateMethods
	Token string `validate:"required"`
}

func (v *validatorCardByToken) validate() httpError {
	return v.vm.structs(v)
}

func (vm *validateMethods) GetCardByToken(ctx *gin.Context) (*validatorCardByToken, httpError) {
	v := &validatorCardByToken{
		vm:    vm,
		Token: ctx.Param("token"),
	}

	return v, v.validate()
}
