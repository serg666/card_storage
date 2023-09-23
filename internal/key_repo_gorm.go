package internal

import (
	"gorm.io/gorm"

	"github.com/serg666/card_storage/models"
)

type gormKeyRepository struct {
	db *gorm.DB
}

type KeySpecificationWithRandom struct{}

func (kswr *KeySpecificationWithRandom) AddGormCondition(db *gorm.DB) *gorm.DB {
	return db.Order("random()").Limit(1)
}

func (r *gormKeyRepository) Create(k *models.Key) (*models.Key, error) {
	m := models.KeyToGormKey(k)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m.ToKey(), nil
}

func (r *gormKeyRepository) Query(specification KeySpecification) ([]*models.Key, error) {
	var gormKeys []models.GormKey

	q := specification.AddGormCondition(r.db).Find(&gormKeys)

	if q.Error != nil {
		return nil, q.Error
	}

	results := make([]*models.Key, 0)
	for _, m := range gormKeys {
		results = append(results, m.ToKey())
	}
	return results, nil
}

func NewKeySpecificationWithRandom() KeySpecification {
	return &KeySpecificationWithRandom{}
}

func NewGormKeyRepository(db *gorm.DB) KeyRepository {
	return &gormKeyRepository{db}
}
