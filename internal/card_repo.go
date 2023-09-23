package internal

import (
	"gorm.io/gorm"

	"github.com/serg666/card_storage/models"
)

type CardSpecification interface {
	AddGormCondition(*gorm.DB) *gorm.DB
}

type CardRepository interface {
	Create(*models.Card) (*models.Card, error)
	Update(*models.Card) error
	ListByPan(*models.Card, *models.Paginator) ([]*models.Card, error)
	GetByToken(string) (*models.Card, error)
	Query(specification CardSpecification) ([]*models.Card, error)
	Delete(*models.Card) error
}
