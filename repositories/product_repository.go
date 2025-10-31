package repositories

import (
	"1kosmetika-marketplace-backend/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll() ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	FindByIDs(ids []uint) ([]models.Product, error)
	// ДОБАВЛЯЕМ новые методы
	FindWithPagination(page, limit int) ([]models.Product, int64, error)
	FindWithFilters(filter ProductFilter, page, limit int) ([]models.Product, int64, error)
	GetCategories() ([]string, error)
	GetBrands() ([]string, error)
}

// ДОБАВЛЯЕМ структуру для фильтров
type ProductFilter struct {
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	Search   string  `json:"search"`
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) FindByIDs(ids []uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("id IN ?", ids).Find(&products).Error
	return products, err
}


func (r *productRepository) FindWithPagination(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * limit
	

	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get products with pagination
	err := r.db.Limit(limit).Offset(offset).Find(&products).Error
	
	return products, total, err
}


func (r *productRepository) FindWithFilters(filter ProductFilter, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})
	
	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}


	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Find(&products).Error
	
	return products, total, err
}


func (r *productRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Product{}).Distinct().Pluck("category", &categories).Error
	return categories, err
}


func (r *productRepository) GetBrands() ([]string, error) {
	var brands []string
	err := r.db.Model(&models.Product{}).Distinct().Pluck("brand", &brands).Error
	return brands, err
}