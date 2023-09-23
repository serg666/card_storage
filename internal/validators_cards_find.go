package internal

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serg666/card_storage/models"
)

type validatorCardsFind struct {
	vm            *validateMethods
	Pan           string  `validate:"required,numeric"`
	ExpDate       *string `validate:"omitempty,notempty"`
	Holder        *string `validate:"omitempty,notempty"`
	Limit         int     `validate:"omitempty,numeric,min=1,max=100"`
	Offset        int     `validate:"omitempty,numeric,min=1"`
	expDateAsTime *time.Time
}

func (v *validatorCardsFind) validateExpDate() error {
	if v.ExpDate == nil {
		return nil
	}

	if len(*v.ExpDate) != 5 {
		return errors.New("invalid date length")
	}

	expDateAsTime, err := time.Parse(models.EXPIRE_DATE_FORMAT, *v.ExpDate)
	if err != nil {
		return err
	}

	v.expDateAsTime = &expDateAsTime
	return nil
}

func (v *validatorCardsFind) validate() httpError {
	if v.validateExpDate() != nil {
		return &httpValidationError{Fields: map[string]string{
			"ExpDate": "Invalid date format",
		}}
	}

	return v.vm.structs(v)
}

func (v *validatorCardsFind) filtersModel() *models.Card {
	return &models.Card{
		PAN:     v.Pan,
		Holder:  v.Holder,
		ExpDate: v.expDateAsTime,
	}
}

func (v *validatorCardsFind) paginator() *models.Paginator {
	return &models.Paginator{
		Limit:  v.Limit,
		Offset: v.Offset,
	}
}

type CardListRequest struct {
	Limit   int     `form:"limit,default=10"`
	Offset  int     `form:"offset,default=0"`
	Pan     string  `form:"pan" validate:"required,numeric"`
	Holder  *string `form:"holder" validate:"omitempty"`
	ExpDate *string `form:"exp_date" validate:"omitempty"`
}

func (vm *validateMethods) CardsList(ctx *gin.Context) (*validatorCardsFind, httpError) {
	var req CardListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, httpErrValidationFailed
	}

	v := &validatorCardsFind{
		vm:      vm,
		Limit:   req.Limit,
		Offset:  req.Offset,
		Pan:     models.CleanPan(req.Pan),
		Holder:  req.Holder,
		ExpDate: req.ExpDate,
	}

	return v, v.validate()
}
