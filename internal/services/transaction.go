package services

import (
	"wallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService interface {
	GetTransactionByRef(ref string) (*models.Transaction, error)
	GetTransactionsByAccountID(accountID uuid.UUID) ([]models.Transaction, error)
}

type transactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) TransactionService {
	return &transactionService{db: db}
}

func (s *transactionService) GetTransactionByRef(ref string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := s.db.First(&transaction, "ref = ?", ref).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *transactionService) GetTransactionsByAccountID(accountID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.db.Where("account_id = ?", accountID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
