package models

import (
	"time"

	"gorm.io/gorm"
)

type WordCategory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CustomerID uint   `json:"customer_id" gorm:"uniqueIndex:idx_customer_word_cat"`
	Name       string `json:"name" gorm:"uniqueIndex:idx_customer_word_cat"`
	Color      int    `json:"color"` // hue 0–360
}
