package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CustomerID uint       `json:"customer_id" gorm:"index"`
	Name       string     `json:"name"`
	Icon       string     `json:"icon"`
	Hue        int        `json:"hue"`
	Sentences  []Sentence `json:"sentences" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

type Sentence struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CategoryID uint   `json:"category_id" gorm:"index"`
	Text       string `json:"text"`
	Meaning    string `json:"meaning"`
	Note       string `json:"note"`
}
