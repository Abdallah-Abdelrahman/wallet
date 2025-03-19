package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user information.
type User struct {
	ID        uuid.UUID `gorm:"type:TEXT;primaryKey"`
	Email     string    `gorm:"unique;not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Accounts  []Account      `gorm:"foreignKey:UserID"`
}

// BeforeCreate generates a new UUID for the ID field.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
