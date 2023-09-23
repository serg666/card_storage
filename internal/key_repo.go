package internal

import (
	"gorm.io/gorm"

	"github.com/serg666/card_storage/models"
)

type KeySpecification interface {
	AddGormCondition(*gorm.DB) *gorm.DB
}

type KeyRepository interface {
	Create(*models.Key) (*models.Key, error)
	Query(specification KeySpecification) ([]*models.Key, error)
}
