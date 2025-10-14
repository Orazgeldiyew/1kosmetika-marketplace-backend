package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type NotificationService interface {
	CreateNotification(notification *models.Notification) error
	GetUserNotifications(userID uint, page, limit int) ([]models.Notification, error)
	GetUnreadCount(userID uint) (int64, error)
	MarkAsRead(notificationID uint, userID uint) error
	MarkAllAsRead(userID uint) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
}

func NewNotificationService(notificationRepo repositories.NotificationRepository, userRepo repositories.UserRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) CreateNotification(notification *models.Notification) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(notification.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetUserNotifications(userID uint, page, limit int) ([]models.Notification, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.notificationRepo.FindByUserID(userID, page, limit)
}

func (s *notificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *notificationService) MarkAsRead(notificationID uint, userID uint) error {
	notification, err := s.notificationRepo.FindByID(notificationID)
	if err != nil {
		return fmt.Errorf("notification not found")
	}

	// Check if notification belongs to user
	if notification.UserID != userID {
		return fmt.Errorf("access denied")
	}

	notification.IsRead = true
	return s.notificationRepo.Update(notification)
}

func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}