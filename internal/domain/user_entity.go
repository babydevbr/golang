// Package domain entities
package domain

import (
	"time"
)

// User entity.
type User struct {
	ID          string     `gorm:"primaryKey;unique;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string     `json:"name" validate:"required"`
	Email       string     `gorm:"uniqueIndex;not null" validate:"required,email" json:"email"`
	Phone       string     `gorm:"index;not null" validate:"required" json:"whatsapp"`
	Password    string     `gorm:"index;" validate:"required" json:"password,omitempty"`
	Role        string     `gorm:"index;default:DEV" json:"role"`
	ActivatedAt *time.Time `gorm:"index" json:"activatedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}
