package models

import (
	"time"
)

type Key struct {
	ID   uint
	Name string
}

type GormKey struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(100);unique_index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (kg GormKey) ToKey() *Key {
	return &Key{
		ID:   kg.ID,
		Name: kg.Name,
	}
}

func KeyToGormKey(k *Key) *GormKey {
	return &GormKey{
		Name: k.Name,
	}
}
