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

	Name     string  `json:"name"`
	Email    *string `json:"email" gorm:"uniqueIndex"` // nullable; existing DB: ALTER TABLE customers ALTER COLUMN email DROP NOT NULL;
	Phone    *string `json:"phone" gorm:"uniqueIndex"` // nullable
	Password string  `json:"-"`
	Image    string  `json:"image"`
	Plan     string  `json:"plan"`   // "pro_monthly", "pro_annual", "free"
	Role     string  `json:"role"`   // "learner", "admin"
	Streak         int    `json:"streak"`
	LastActiveDate string `json:"last_active_date"` // YYYY-MM-DD
	Status         string `json:"status"`           // "active", "inactive"
}

// CustomerIdentity links an OAuth provider or phone to one Customer.
// One customer can have many identities (Google + Facebook + phone).
type CustomerIdentity struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	CustomerID uint   `json:"customer_id" gorm:"index"`
	Provider   string `json:"provider" gorm:"uniqueIndex:idx_identity"`    // "google", "facebook", "phone"
	ProviderID string `json:"provider_id" gorm:"uniqueIndex:idx_identity"` // sub / fb-id / phone
}

// OTPCode stores phone OTP challenges. Short-lived; not soft-deleted.
type OTPCode struct {
	ID        uint      `gorm:"primaryKey"`
	Phone     string    `gorm:"index"`
	Code      string
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
}
