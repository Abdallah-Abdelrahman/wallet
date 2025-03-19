package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Define TransactionType as a custom string type
type TransactionType string

// Define constants for each type of transaction
const (
	TopUp  TransactionType = "top-up"
	Charge TransactionType = "charge"
)

// Transaction represents a account transactions.
type Transaction struct {
	ID              uuid.UUID       `gorm:"type:TEXT;primaryKey"`
	TransactionType TransactionType `gorm:"type:varchar(10);not null;check:transaction_type IN ('top-up', 'charge')"`
	Amount          float64         `gorm:"type:decimal(10,2);not null"`
	Ref             string          `gorm:"not null;unique"`
	AccountID       uuid.UUID       `gorm:"type:uuid;not null"`
	Account         Account         `gorm:"foreignKey:AccountID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate generates a new UUID for the ID field.
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
