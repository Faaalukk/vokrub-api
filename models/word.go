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
	CategoryID *uint       `json:"category_id"`
	Pos        StringSlice `json:"pos" gorm:"type:text;default:'[]'"`
	Translate  string      `json:"translate"` // native-language translation (e.g. Thai)
	Meaning    string      `json:"meaning"`   // English definition
	Note       string      `json:"note"`
	Synonyms   StringSlice `json:"synonyms" gorm:"type:text;default:'[]'"`
	Box        int    `json:"box" gorm:"default:1"`
	Seen       int    `json:"seen" gorm:"default:0"`
	Due        bool   `json:"due" gorm:"default:true"`
	Added      string `json:"added"` // YYYY-MM-DD
}
