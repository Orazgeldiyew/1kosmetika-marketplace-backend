package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	FindCartItem(cartID, productID uint) (*models.CartItem, error)
	FindCartItemByID(itemID uint) (*models.CartItem, error) // ДОБАВИТЬ этот метод
	CreateCartItem(item *models.CartItem) error
	UpdateCartItem(item *models.CartItem) error
	DeleteCartItem(itemID uint) error
	ClearCart(cartID uint) error
	GetCartWithItems(userID uint) (*models.Cart, error)
	DeleteCartItemByProduct(cartID, productID uint) error // ДОБАВИТЬ этот метод
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepository) FindCartItem(cartID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ДОБАВИТЬ этот метод
func (r *cartRepository) FindCartItemByID(itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Preload("Product").First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) CreateCartItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateCartItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) DeleteCartItem(itemID uint) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}

// ДОБАВИТЬ этот метод
func (r *cartRepository) DeleteCartItemByProduct(cartID, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCartWithItems(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}