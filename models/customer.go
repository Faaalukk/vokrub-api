package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"-"`
	Image    string `json:"image"`
	Plan     string `json:"plan"`   // "pro_monthly", "pro_annual", "free"
	Role     string `json:"role"`   // "learner", "admin"
	Streak   int    `json:"streak"`
	Status   string `json:"status"` // "active", "inactive"
}
