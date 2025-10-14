package services

import (
	"fmt"
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(id uint, product *models.Product) error
	DeleteProduct(id uint) error
	GetProductsByIDs(ids []uint) ([]models.Product, error)
	// ДОБАВЛЯЕМ новые методы в интерфейс
	GetProductsWithPagination(page, limit int) ([]models.Product, int64, error)
	GetProductsWithFilters(filter repositories.ProductFilter, page, limit int) ([]models.Product, int64, error)
	GetCategories() ([]string, error)
	GetBrands() ([]string, error)
	ValidateProduct(product *models.Product) error
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(product *models.Product) error {
	if err := s.ValidateProduct(product); err != nil {
		return err
	}
	return s.productRepo.Create(product)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

func (s *productService) UpdateProduct(id uint, product *models.Product) error {
	if err := s.ValidateProduct(product); err != nil {
		return err
	}
	
	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Update fields
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.ImageURL = product.ImageURL
	existingProduct.Category = product.Category
	existingProduct.Brand = product.Brand
	existingProduct.Stock = product.Stock

	return s.productRepo.Update(existingProduct)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *productService) GetProductsByIDs(ids []uint) ([]models.Product, error) {
	return s.productRepo.FindByIDs(ids)
}

// ДОБАВЛЯЕМ РЕАЛИЗАЦИЮ новых методов:

func (s *productService) GetProductsWithPagination(page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	return s.productRepo.FindWithPagination(page, limit)
}

func (s *productService) GetProductsWithFilters(filter repositories.ProductFilter, page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	return s.productRepo.FindWithFilters(filter, page, limit)
}

func (s *productService) GetCategories() ([]string, error) {
	return s.productRepo.GetCategories()
}

func (s *productService) GetBrands() ([]string, error) {
	return s.productRepo.GetBrands()
}

func (s *productService) ValidateProduct(product *models.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if len(product.Name) < 2 {
		return fmt.Errorf("product name must be at least 2 characters")
	}
	if product.Price <= 0 {
		return fmt.Errorf("product price must be positive")
	}
	if product.Category == "" {
		return fmt.Errorf("product category is required")
	}
	if product.Brand == "" {
		return fmt.Errorf("product brand is required")
	}
	if product.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}
	return nil
}