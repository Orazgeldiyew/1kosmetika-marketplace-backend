// services/cart_service.go
package services

import (
	"fmt"

	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type CartService interface {
	GetCart(userID uint) (*models.Cart, error)
	AddToCart(userID uint, productID uint, quantity int) error
	UpdateCartItem(userID uint, itemID uint, quantity int) error
	RemoveFromCart(userID uint, itemID uint) error
	ClearCart(userID uint) error
}

type cartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// Возвращаем корзину; если не найдена — "пустая" структура (ID=0, Items=[]), без ошибки.
func (s *cartService) GetCart(userID uint) (*models.Cart, error) {
	cart, err := s.cartRepo.GetCartWithItems(userID)
	if err != nil {
		return &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
		}, nil
	}
	return cart, nil
}

func (s *cartService) AddToCart(userID uint, productID uint, quantity int) error {
	if quantity <= 0 {
		quantity = 1
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return fmt.Errorf("product not found")
	}
	if product.Stock < quantity {
		return fmt.Errorf("not enough stock available")
	}

	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		cart = &models.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return fmt.Errorf("failed to create cart")
		}
	}

	existingItem, err := s.cartRepo.FindCartItem(cart.ID, productID)
	if err == nil {
		newQuantity := existingItem.Quantity + quantity
		if product.Stock < newQuantity {
			return fmt.Errorf("not enough stock available")
		}
		existingItem.Quantity = newQuantity
		existingItem.Price = product.Price * float64(newQuantity)
		return s.cartRepo.UpdateCartItem(existingItem)
	}

	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price * float64(quantity),
	}
	return s.cartRepo.CreateCartItem(cartItem)
}

func (s *cartService) UpdateCartItem(userID uint, itemID uint, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be >= 1")
	}

	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return fmt.Errorf("cart not found")
	}

	cartItem, err := s.cartRepo.FindCartItemByID(itemID)
	if err != nil {
		return fmt.Errorf("cart item not found")
	}
	if cartItem.CartID != cart.ID {
		return fmt.Errorf("cart item not found")
	}

	product, err := s.productRepo.FindByID(cartItem.ProductID)
	if err != nil {
		return fmt.Errorf("product not found")
	}
	if product.Stock < quantity {
		return fmt.Errorf("not enough stock available")
	}

	cartItem.Quantity = quantity
	cartItem.Price = product.Price * float64(quantity)
	return s.cartRepo.UpdateCartItem(cartItem)
}

// Идемпотентное удаление: отсутствие/чужой item — не ошибка.
func (s *cartService) RemoveFromCart(userID uint, itemID uint) error {
	_, err := s.cartRepo.DeleteCartItemOwnedByUser(userID, itemID)
	return err
}

// Идемпотентная очистка: если корзины нет — считаем, что уже пусто.
func (s *cartService) ClearCart(userID uint) error {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil
	}
	return s.cartRepo.ClearCart(cart.ID)
}
