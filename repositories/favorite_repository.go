package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(favorite *models.Favorite) error
	Delete(userID, productID uint) error
	FindByUserAndProduct(userID, productID uint) (*models.Favorite, error)
	FindByUserID(userID uint) ([]models.Favorite, error)
	Exists(userID, productID uint) (bool, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Create(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepository) Delete(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Favorite{}).Error
}

func (r *favoriteRepository) FindByUserAndProduct(userID, productID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) FindByUserID(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.Preload("Product").Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepository) Exists(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}