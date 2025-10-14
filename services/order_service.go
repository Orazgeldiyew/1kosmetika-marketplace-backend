package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type OrderService interface {
	CreateOrder(userID uint, productIDs []uint) (*models.Order, error)
	GetUserOrders(userID uint) ([]models.Order, error)
	GetOrderByID(id uint) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

type orderService struct {
	orderRepo        repositories.OrderRepository
	productRepo      repositories.ProductRepository
	cartRepo         repositories.CartRepository
	notificationRepo repositories.NotificationRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
	notificationRepo repositories.NotificationRepository,
) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		productRepo:      productRepo,
		cartRepo:         cartRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *orderService) CreateOrder(userID uint, productIDs []uint) (*models.Order, error) {
	// Get products
	products, err := s.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if len(products) != len(productIDs) {
		return nil, fmt.Errorf("some products not found")
	}

	// Calculate total
	var total float64
	for _, product := range products {
		total += product.Price
	}

	// Create order
	order := &models.Order{
		UserID:   userID,
		Products: products,
		Total:    total,
		Status:   "pending",
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create notification (как в старом коде)
	notification := &models.Notification{
		UserID:  userID,
		Title:   "Заказ оформлен",
		Message: fmt.Sprintf("Ваш заказ #%d успешно оформлен. Сумма: %.2f руб.", order.ID, total),
		Type:    "success",
	}
	
	if err := s.notificationRepo.Create(notification); err != nil {
		// Логируем ошибку, но не прерываем создание заказа
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	// Clear user's cart
	cart, err := s.cartRepo.FindByUserID(userID)
	if err == nil && cart != nil {
		s.cartRepo.ClearCart(cart.ID)
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *orderService) GetOrderByID(id uint) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *orderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}