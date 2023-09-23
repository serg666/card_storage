package internal

import (
	"time"

	"gorm.io/gorm"

	vault "github.com/mittwald/vaultgo"
	"github.com/serg666/card_storage/models"
)

type gormSessionRepository struct {
	db    *gorm.DB
	vault *vault.Client
}

type SessionSpecificationWithoutKey struct{}

func (sswok *SessionSpecificationWithoutKey) AddGormCondition(db *gorm.DB) *gorm.DB {
	return db.Where("c_key_id is null")
}

func (r *gormSessionRepository) Update(s *models.Session) error {
	m, err := models.SessionToGormSession(s, r.vault)

	if err != nil {
		return err
	}

	m.ID = s.ID

	return r.db.Model(&m).Updates(m).Error
}

func (r *gormSessionRepository) Query(specification SessionSpecification) ([]*models.Session, error) {
	var gormSessions []models.GormSession

	q := specification.AddGormCondition(r.db.Joins("CKey")).Find(&gormSessions)

	if q.Error != nil {
		return nil, q.Error
	}

	results := make([]*models.Session, 0)
	for _, m := range gormSessions {
		session, err := m.ToSession(r.vault)

		if err != nil {
			return nil, err
		}

		results = append(results, session)
	}
	return results, nil
}

func (r *gormSessionRepository) Create(s *models.Session) (*models.Session, error) {
	m, err := models.SessionToGormSession(s, r.vault)
	if err != nil {
		return nil, err
	}

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	return m.ToSession(r.vault)
}

func (r *gormSessionRepository) Get(key string) (*models.Session, error) {
	m := &models.GormSession{}
	tx := r.db.Joins("CKey").
		Where("key = ? AND expire_at > ?", key, time.Now()).
		Order("id DESC").
		First(&m)
	if err := tx.Error; err != nil {
		return nil, err
	}

	return m.ToSession(r.vault)
}

func (r *gormSessionRepository) DeleteExpired() error {
	return r.db.Where("expire_at <= ?", time.Now()).Delete(&models.GormSession{}).Error
}

func NewSessionSpecificationWithoutKey() SessionSpecification {
	return &SessionSpecificationWithoutKey{}
}

func NewGormSessionRepository(db *gorm.DB, vault *vault.Client) SessionRepository {
	return &gormSessionRepository{db, vault}
}
