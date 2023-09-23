package internal

import (
	"github.com/stretchr/testify/mock"

	"github.com/serg666/card_storage/models"
)

type mockCardRepository struct {
	mock.Mock
}

func (m *mockCardRepository) Create(c *models.Card) (*models.Card, error) {
	args := m.Called(c)
	return args.Get(0).(*models.Card), args.Error(1)
}

func (m *mockCardRepository) Update(c *models.Card) error {
	args := m.Called(c)
	return args.Error(1)
}

func (m *mockCardRepository) ListByPan(c *models.Card, p *models.Paginator) ([]*models.Card, error) {
	args := m.Called(c)
	return args.Get(0).([]*models.Card), args.Error(1)
}

func (m *mockCardRepository) GetByToken(token string) (*models.Card, error) {
	args := m.Called(token)
	return args.Get(0).(*models.Card), args.Error(1)
}

func (m *mockCardRepository) Query(specification CardSpecification) ([]*models.Card, error) {
	args := m.Called(specification)
	return args.Get(0).([]*models.Card), args.Error(1)
}

func (m *mockCardRepository) Delete(c *models.Card) error {
	args := m.Called(c)
	return args.Error(0)
}

func NewMockCardRepository() *mockCardRepository {
	return &mockCardRepository{}
}
