// repositories/cart_repository.go
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
	FindCartItemByID(itemID uint) (*models.CartItem, error)
	CreateCartItem(item *models.CartItem) error
	UpdateCartItem(item *models.CartItem) error
	DeleteCartItem(itemID uint) error
	DeleteCartItemByProduct(cartID, productID uint) error
	DeleteCartItemOwnedByUser(userID, itemID uint) (bool, error)

	ClearCart(cartID uint) error
	GetCartWithItems(userID uint) (*models.Cart, error)
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

func (r *cartRepository) DeleteCartItemByProduct(cartID, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&models.CartItem{}).Error
}

func (r *cartRepository) DeleteCartItemOwnedByUser(userID, itemID uint) (bool, error) {
	res := r.db.
		Where(`id = ? AND cart_id IN (SELECT id FROM carts WHERE user_id = ?)`, itemID, userID).
		Delete(&models.CartItem{})

	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
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
