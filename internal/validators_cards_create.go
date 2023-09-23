package internal

import (
	"errors"
	"io"
	"time"

	"github.com/serg666/card_storage/models"
)

type validatorCardsCreate struct {
	vm            *validateMethods
	Pan           string  `json:"pan" validate:"required"`
	Holder        *string `json:"holder" validate:"omitempty,notempty"`
	ExpDate       *string `json:"exp_date" validate:"omitempty,notempty"`
	expDateAsTime *time.Time
}

func (v *validatorCardsCreate) validate(src io.Reader) httpError {
	if err := v.vm.json(v, src); err != nil {
		return err
	}

	v.Pan = models.CleanPan(v.Pan)
	if err := v.vm.structs(v); err != nil {
		return err
	}

	if v.validateExpDate() != nil {
		return &httpValidationError{Fields: map[string]string{
			"ExpDate": "Invalid date format",
		}}
	}

	return v.vm.structs(v)
}

func (v *validatorCardsCreate) validateExpDate() error {
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

func (v *validatorCardsCreate) toCardModel(key *models.Key) *models.Card {
	return &models.Card{
		PAN:     v.Pan,
		Holder:  v.Holder,
		ExpDate: v.expDateAsTime,
		Key:     key,
	}
}

func (vm *validateMethods) CardsCreate(src io.Reader) (*validatorCardsCreate, httpError) {
	v := &validatorCardsCreate{
		vm: vm,
	}
	return v, v.validate(src)
}
