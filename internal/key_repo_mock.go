package internal

import (
	"github.com/stretchr/testify/mock"

	"github.com/serg666/card_storage/models"
)

type mockKeyRepository struct {
	mock.Mock
}

func (m *mockKeyRepository) Create(k *models.Key) (*models.Key, error) {
	args := m.Called(k)
	return args.Get(0).(*models.Key), args.Error(1)
}

func (m *mockKeyRepository) Query(specification KeySpecification) ([]*models.Key, error) {
	args := m.Called(specification)
	key := &models.Key{ID: 1, Name: "1"}
	list := make([]*models.Key, 0)
	list = append(list, key)
	return list, args.Error(1)
}

func NewMockKeyRepository() *mockKeyRepository {
	return &mockKeyRepository{}
}
