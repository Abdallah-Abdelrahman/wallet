package services

import (
	"errors"
	"fmt"
	"math"
	"time"
	"wallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccountWithUser(email, firstName, lastName string) (*models.Account, error)
	GetAccountByID(accountID uuid.UUID) (*models.Account, error)
	TopUp(accountID uuid.UUID, amount float64) (*models.Transaction, error)
	Charge(accountID uuid.UUID, amount float64) (*models.Transaction, error)
}

type accountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{db: db}
}

// CreateAccountWithUser creates a new user and a corresponding account with a 0.00 balance.
func (s *accountService) CreateAccountWithUser(email, firstName, lastName string) (*models.Account, error) {
	userService := NewUserService(s.db)

	// Check if user already exists
	_, user_err := userService.GetUserByEmail(email)
	if user_err == nil {
		return nil, errors.New("user already exists")
	}

	// Start a transaction
	tx := s.db.Begin()

	// Create the user
	new_user := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	// Save the new user
	if err := tx.Create(&new_user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create the account for the user with an initial balance of 0.00
	account := &models.Account{
		Balance: 0.00,
		UserID:  new_user.ID,
	}

	// Save the new account
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// set the relationship
	account.User = *new_user

	return account, nil
}

// GetAccountByID retrieves an account by its ID.
func (s *accountService) GetAccountByID(accountID uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := s.db.First(&account, "id = ?", accountID).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// TopUp adds funds to an account.
func (s *accountService) TopUp(accountID uuid.UUID, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Retrieve the account
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Add balance to the account
	account.Balance += amount
	account.Balance = math.Round(account.Balance*100) / 100
	// Retrieve the user associated with the account
	acc_u, u_err := NewUserService(s.db).GetUserByID(account.UserID)

	if u_err != nil {
		tx.Rollback()
		return nil, u_err
	}

	var existingTransaction models.Transaction

	// Generate a unique Ref for each transaction
	ref := fmt.Sprintf("TXN-%s-%d", uuid.New().String(), time.Now().UnixNano())

	// Check if a transaction with the same Ref already exists
	if err := s.db.Where("ref = ?", ref).First(&existingTransaction).Error; err == nil {
		return nil, errors.New("duplicate transaction detected")
	}

	// Create a new transaction
	transaction := &models.Transaction{
		TransactionType: models.TopUp,
		Amount:          math.Round(amount*100) / 100,
		Ref:             ref,
		AccountID:       accountID,
	}

	// Save the transaction
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update the account with the new balance
	if err := tx.Save(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	// set the relationship
	transaction.Account = *account
	transaction.Account.User = *acc_u

	return transaction, nil
}

// Charge deducts funds from an account.
func (s *accountService) Charge(accountID uuid.UUID, amount float64) (*models.Transaction, error) {
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

	var existingTransaction models.Transaction

	// Generate a unique Ref for each transaction
	ref := fmt.Sprintf("TXN-%s-%d", uuid.New().String(), time.Now().UnixNano())

	// Check if a transaction with the same Ref already exists
	if err := s.db.Where("ref = ?", ref).First(&existingTransaction).Error; err == nil {
		return nil, errors.New("duplicate transaction detected")
	}

	transaction := &models.Transaction{
		TransactionType: models.Charge,
		Amount:          math.Round(amount*100) / 100,
		Ref:             ref,
		AccountID:       accountID,
	}

	// Deduct balance from the account
	account.Balance -= amount
	account.Balance = math.Round(account.Balance*100) / 100

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

	// Retrieve the user associated with the account
	acc_u, u_err := NewUserService(s.db).GetUserByID(account.UserID)

	if u_err != nil {
		tx.Rollback()
		return nil, u_err
	}

	// set the relationship
	transaction.Account = *account
	transaction.Account.User = *acc_u
	return transaction, nil
}
