package services

import (
	"wallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
}

type userService struct {
	db *gorm.DB
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}
