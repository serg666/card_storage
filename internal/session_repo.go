package internal

import (
	"gorm.io/gorm"

	"github.com/serg666/card_storage/models"
)

type SessionSpecification interface {
	AddGormCondition(*gorm.DB) *gorm.DB
}

type SessionRepository interface {
	Create(*models.Session) (*models.Session, error)
	Update(*models.Session) error
	Get(key string) (*models.Session, error)
	Query(specification SessionSpecification) ([]*models.Session, error)
	DeleteExpired() error
}
