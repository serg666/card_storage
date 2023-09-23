package internal

import (
	"gorm.io/gorm"

	vault "github.com/mittwald/vaultgo"
	"github.com/serg666/card_storage/models"
)

type gormCardRepository struct {
	db    *gorm.DB
	vault *vault.Client
}

type CardSpecificationWithoutKey struct{}

func (cswok *CardSpecificationWithoutKey) AddGormCondition(db *gorm.DB) *gorm.DB {
	return db.Where("key_id is null")
}

func (r *gormCardRepository) Create(c *models.Card) (*models.Card, error) {
	if err := c.GenerateToken(); err != nil {
		return nil, err
	}

	m, err := models.CardToGormCard(c, r.vault)
	if err != nil {
		return nil, err
	}

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m.ToCard(r.vault)
}

func (r *gormCardRepository) Update(c *models.Card) error {
	m, err := models.CardToGormCard(c, r.vault)

	if err != nil {
		return err
	}

	m.ID = c.ID
	m.MASK = c.PanMask()

	return r.db.Model(&m).Updates(m).Error
}

func (r *gormCardRepository) Query(specification CardSpecification) ([]*models.Card, error) {
	var gormCards []models.GormCard

	q := specification.AddGormCondition(r.db.Joins("Key")).Find(&gormCards)

	if q.Error != nil {
		return nil, q.Error
	}

	results := make([]*models.Card, 0)
	for _, m := range gormCards {
		card, err := m.ToCard(r.vault)

		if err != nil {
			return nil, err
		}

		results = append(results, card)
	}
	return results, nil
}

func (r *gormCardRepository) ListByPan(c *models.Card, p *models.Paginator) ([]*models.Card, error) {
	var gormCards []models.GormCard

	q := r.db.Joins("Key").Where(&models.GormCard{MASK: c.PanMask(), Holder: c.Holder, ExpDate: c.ExpDate})

	if c.Holder == nil {
		q = q.Where("holder is NULL")
	}

	if c.ExpDate == nil {
		q = q.Where("exp_date is NULL")
	}

	q = q.Order("id DESC").Find(&gormCards)

	if q.Error != nil {
		return nil, q.Error
	}

	results := make([]*models.Card, 0)
	i := 0
	for _, m := range gormCards {
		card, err := m.ToCard(r.vault)

		if err != nil {
			return nil, err
		}

		if c.PAN == card.PAN && i >= p.Offset && i < p.Offset+p.Limit {
			results = append(results, card)
			i++
		}
	}

	return results, nil
}

func (r *gormCardRepository) GetByToken(token string) (*models.Card, error) {
	var gormCard models.GormCard

	q := r.db.Joins("Key").Where(&models.GormCard{
		Token: token,
	}).Order("id DESC").First(&gormCard)
	if q.Error != nil {
		return nil, q.Error
	}

	return gormCard.ToCard(r.vault)
}

func (r *gormCardRepository) Delete(c *models.Card) error {
	return r.db.Delete(&models.GormCard{}, c.ID).Error
}

func NewCardSpecificationWithoutKey() CardSpecification {
	return &CardSpecificationWithoutKey{}
}

func NewGormCardRepository(db *gorm.DB, vault *vault.Client) CardRepository {
	return &gormCardRepository{db, vault}
}
