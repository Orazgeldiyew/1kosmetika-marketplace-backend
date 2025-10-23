package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	CreateOrderProducts(items []models.OrderProduct) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderProducts(items []models.OrderProduct) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Products").Preload("User").First(&order, id).Error
	return &order, err
}

func (r *orderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Products").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Products").Preload("User").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
