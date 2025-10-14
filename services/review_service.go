package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type ReviewService interface {
	CreateReview(userID uint, productID uint, rating int, comment string) (*models.Review, error)
	UpdateReview(reviewID uint, userID uint, rating int, comment string) (*models.Review, error)
	DeleteReview(reviewID uint, userID uint) error
	GetProductReviews(productID uint) ([]models.Review, float64, int, error)
	GetUserReviews(userID uint) ([]models.Review, error)
	GetReviewByID(reviewID uint) (*models.Review, error)
}

type reviewService struct {
	reviewRepo  repositories.ReviewRepository
	productRepo repositories.ProductRepository
}

func NewReviewService(reviewRepo repositories.ReviewRepository, productRepo repositories.ProductRepository) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
	}
}

func (s *reviewService) CreateReview(userID uint, productID uint, rating int, comment string) (*models.Review, error) {
	// Check if product exists
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if user already reviewed this product
	existingReview, err := s.reviewRepo.FindByUserAndProduct(userID, productID)
	if err == nil && existingReview != nil {
		return nil, fmt.Errorf("you have already reviewed this product")
	}

	review := &models.Review{
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Get the created review with relations
	return s.reviewRepo.FindByID(review.ID)
}

func (s *reviewService) UpdateReview(reviewID uint, userID uint, rating int, comment string) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("review not found")
	}

	// Check if user owns the review
	if review.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	review.Rating = rating
	review.Comment = comment

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return s.reviewRepo.FindByID(review.ID)
}

func (s *reviewService) DeleteReview(reviewID uint, userID uint) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}

	// Check if user owns the review
	if review.UserID != userID {
		return fmt.Errorf("access denied")
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *reviewService) GetProductReviews(productID uint) ([]models.Review, float64, int, error) {
	reviews, err := s.reviewRepo.FindByProductID(productID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get reviews: %w", err)
	}

	avgRating, totalCount, err := s.reviewRepo.GetProductStats(productID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get review stats: %w", err)
	}

	return reviews, avgRating, totalCount, nil
}

func (s *reviewService) GetUserReviews(userID uint) ([]models.Review, error) {
	return s.reviewRepo.FindByUserID(userID)
}

func (s *reviewService) GetReviewByID(reviewID uint) (*models.Review, error) {
	return s.reviewRepo.FindByID(reviewID)
}