package services

import (
	"errors"
	"wallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(userID uuid.UUID) (*models.Account, error)
	GetAccountByID(accountID uuid.UUID) (*models.Account, error)
	TopUp(accountID uuid.UUID, amount float64, ref string) (*models.Transaction, error)
	Charge(accountID uuid.UUID, amount float64, ref string) (*models.Transaction, error)
}

type accountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{db: db}
}

func (s *accountService) CreateAccount(userID uuid.UUID) (*models.Account, error) {
	account := &models.Account{
		UserID:  userID,
		Balance: 0.00,
	}
	err := s.db.Create(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByID(accountID uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := s.db.First(&account, "id = ?", accountID).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *accountService) TopUp(accountID uuid.UUID, amount float64, ref string) (*models.Transaction, error) {
	// Retrieve the account
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := &models.Transaction{
		TransactionType: models.TopUp,
		Amount:          amount,
		Ref:             ref,
		AccountID:       accountID,
	}

	// Add balance to the account
	account.Balance += amount

	// Save the transaction and update the account
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	tx.Commit()

	return transaction, nil
}

func (s *accountService) Charge(accountID uuid.UUID, amount float64, ref string) (*models.Transaction, error) {
	// Retrieve the account
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Ensure there are sufficient funds
	if account.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Create transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction := &models.Transaction{
		TransactionType: models.Charge,
		Amount:          amount,
		Ref:             ref,
		AccountID:       accountID,
	}

	// Deduct balance from the account
	account.Balance -= amount

	// Save the transaction and update the account
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	tx.Commit()

	return transaction, nil
}
