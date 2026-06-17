package models

import (
	"time"

	"gorm.io/gorm"
)

type WordFamily struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `json:"-" gorm:"index"`
	CustomerID uint              `json:"customer_id" gorm:"index"`
	Name       string            `json:"name"`
	Members    []WordFamilyMember `json:"members" gorm:"foreignKey:FamilyID"`
}

type WordFamilyMember struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	FamilyID uint `json:"family_id" gorm:"index;uniqueIndex:idx_family_word"`
	WordID   uint `json:"word_id"   gorm:"uniqueIndex:idx_family_word"`
}
