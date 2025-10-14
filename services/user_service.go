package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
	"1kosmetika-marketplace-backend/utils"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, string, error)
	GetProfile(userID uint) (*models.User, error)
	UpdateRole(userID uint, role string) error
	GetAllUsers() ([]models.User, error)
	DeleteUser(userID uint) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *models.User) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	user.Password = hashedPassword
	user.Role = "user" // default role

	return s.userRepo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return user, token, nil
}

func (s *userService) GetProfile(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateRole(userID uint, role string) error {
	return s.userRepo.UpdateRole(userID, role)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}