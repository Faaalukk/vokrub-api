package models

import (
	"time"

	"gorm.io/gorm"
)

type Word struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CustomerID uint   `json:"customer_id" gorm:"uniqueIndex:idx_customer_word"`
	Word       string `json:"word" gorm:"uniqueIndex:idx_customer_word"`
	CategoryID *uint  `json:"category_id"`
	Pos        string `json:"pos"`
	Meaning    string `json:"meaning"`
	Note       string `json:"note"`
	Box        int    `json:"box" gorm:"default:1"`
	Seen       int    `json:"seen" gorm:"default:0"`
	Due        bool   `json:"due" gorm:"default:true"`
	Added      string `json:"added"` // YYYY-MM-DD
}
