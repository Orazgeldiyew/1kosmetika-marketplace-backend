package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	Update(notification *models.Notification) error
	FindByID(notificationID uint) (*models.Notification, error)
	FindByUserID(userID uint, page, limit int) ([]models.Notification, error)
	GetUnreadCount(userID uint) (int64, error)
	MarkAllAsRead(userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) FindByID(notificationID uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, notificationID).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) FindByUserID(userID uint, page, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	offset := (page - 1) * limit
	
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error
		
	return notifications, err
}

func (r *notificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}