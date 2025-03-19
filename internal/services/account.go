package services

import (
	"errors"
	"time"
	"wallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccountWithUser(email, firstName, lastName string) (*models.Account, error)
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

// CreateAccountWithUser creates a new user and a corresponding account with a 0.00 balance.
func (s *accountService) CreateAccountWithUser(email, firstName, lastName string) (*models.Account, error) {
	// Start a transaction
	tx := s.db.Begin()

	// Create the user
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create the account for the user with an initial balance of 0.00
	account := &models.Account{
		ID:        uuid.New(),
		Balance:   0.00,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
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
