package internal

import (
	"github.com/stretchr/testify/mock"

	"github.com/serg666/card_storage/models"
)

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) Create(s *models.Session) (*models.Session, error) {
	args := m.Called(s)
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *mockSessionRepository) Update(c *models.Session) error {
	args := m.Called(c)
	return args.Error(1)
}

func (m *mockSessionRepository) Query(specification SessionSpecification) ([]*models.Session, error) {
	args := m.Called(specification)
	return args.Get(0).([]*models.Session), args.Error(1)
}

func (m *mockSessionRepository) Get(key string) (*models.Session, error) {
	args := m.Called(key)
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *mockSessionRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func NewMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{}
}
