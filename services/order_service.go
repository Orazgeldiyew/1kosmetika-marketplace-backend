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
	// Load products
	products, err := s.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	if len(products) != len(productIDs) {
		return nil, fmt.Errorf("some products not found")
	}

	// Compute totals & build order items (qty=1 basic case)
	var total float64
	items := make([]models.OrderProduct, 0, len(products))
	for _, p := range products {
		total += p.Price
		items = append(items, models.OrderProduct{
			ProductID: p.ID,
			Quantity:  1,
			Price:     p.Price,
		})
	}

	order := &models.Order{
		UserID:  userID,
		Total:   total,
		Status:  "pending",
		// Keep M2M for quick reads; GORM will maintain join, but we also store richer OrderProduct rows
		Products: products,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Attach order id to items and persist
	for i := range items {
		items[i].OrderID = order.ID
	}
	if err := s.orderRepo.CreateOrderProducts(items); err != nil {
		return nil, fmt.Errorf("failed to save order items: %w", err)
	}

	// Notify
	notification := &models.Notification{
		UserID:  userID,
		Title:   "Заказ оформлен",
		Message: fmt.Sprintf("Ваш заказ #%d успешно оформлен. Сумма: %.2f руб.", order.ID, total),
		Type:    "success",
	}
	_ = s.notificationRepo.Create(notification) // non-blocking

	// Clear cart if exists
	if cart, err := s.cartRepo.FindByUserID(userID); err == nil && cart != nil {
		_ = s.cartRepo.ClearCart(cart.ID)
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
