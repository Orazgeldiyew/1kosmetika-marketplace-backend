package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type FavoriteService interface {
	GetUserFavorites(userID uint) ([]models.Favorite, error)
	AddFavorite(userID, productID uint) error
	RemoveFavorite(userID, productID uint) error
	IsFavorite(userID, productID uint) (bool, error)
}

type favoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	productRepo  repositories.ProductRepository
}

func NewFavoriteService(favoriteRepo repositories.FavoriteRepository, productRepo repositories.ProductRepository) FavoriteService {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		productRepo:  productRepo,
	}
}

func (s *favoriteService) GetUserFavorites(userID uint) ([]models.Favorite, error) {
	return s.favoriteRepo.FindByUserID(userID)
}

func (s *favoriteService) AddFavorite(userID, productID uint) error {

	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found")
	}


	exists, err := s.favoriteRepo.Exists(userID, productID)
	if err != nil {
		return fmt.Errorf("failed to check favorites")
	}
	if exists {
		return fmt.Errorf("product already in favorites")
	}

	favorite := &models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	return s.favoriteRepo.Create(favorite)
}

func (s *favoriteService) RemoveFavorite(userID, productID uint) error {

	exists, err := s.favoriteRepo.Exists(userID, productID)
	if err != nil {
		return fmt.Errorf("failed to check favorites")
	}
	if !exists {
		return fmt.Errorf("product not in favorites")
	}

	return s.favoriteRepo.Delete(userID, productID)
}

func (s *favoriteService) IsFavorite(userID, productID uint) (bool, error) {
	return s.favoriteRepo.Exists(userID, productID)
}