package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	Update(review *models.Review) error
	Delete(reviewID uint) error
	FindByID(reviewID uint) (*models.Review, error)
	FindByUserAndProduct(userID, productID uint) (*models.Review, error)
	FindByProductID(productID uint) ([]models.Review, error)
	FindByUserID(userID uint) ([]models.Review, error)
	GetProductStats(productID uint) (float64, int, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(reviewID uint) error {
	return r.db.Delete(&models.Review{}, reviewID).Error
}

func (r *reviewRepository) FindByID(reviewID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").Preload("Product").First(&review, reviewID).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByUserAndProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByProductID(productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("User").Where("product_id = ?", productID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) FindByUserID(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("Product").Where("user_id = ?", userID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetProductStats(productID uint) (float64, int, error) {
	var result struct {
		AverageRating float64
		TotalCount    int
	}
	
	err := r.db.Model(&models.Review{}).
		Select("AVG(rating) as average_rating, COUNT(*) as total_count").
		Where("product_id = ?", productID).
		Scan(&result).Error

	return result.AverageRating, result.TotalCount, err
}