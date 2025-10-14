package handlers

import (
	"1kosmetika-marketplace-backend/models"
	"1kosmetika-marketplace-backend/repositories"
	"1kosmetika-marketplace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// @Summary Получить все продукты
// @Description Получить список всех продуктов
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Router /api/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary Получить продукты с пагинацией
// @Description Получить продукты с пагинацией
// @Tags products
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(20)
// @Success 200 {object} gin.H
// @Router /api/products/paginated [get]
func (h *ProductHandler) GetProductsPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productService.GetProductsWithPagination(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

// @Summary Поиск и фильтрация продуктов
// @Description Поиск продуктов с фильтрацией и пагинацией
// @Tags products
// @Produce json
// @Param category query string false "Категория"
// @Param brand query string false "Бренд"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param search query string false "Поисковый запрос"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на странице" default(20)
// @Success 200 {object} gin.H
// @Router /api/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	// Parse query parameters
	filter := repositories.ProductFilter{
		Category: c.Query("category"),
		Brand:    c.Query("brand"),
		Search:   c.Query("search"),
	}

	// Parse price filters
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = minPrice
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = maxPrice
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productService.GetProductsWithFilters(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"pages":     (total + int64(limit) - 1) / int64(limit),
		"filters":   filter,
	})
}

// @Summary Получить все категории
// @Description Получить список всех категорий продуктов
// @Tags products
// @Produce json
// @Success 200 {array} string
// @Router /api/products/categories [get]
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// @Summary Получить все бренды
// @Description Получить список всех брендов продуктов
// @Tags products
// @Produce json
// @Success 200 {array} string
// @Router /api/products/brands [get]
func (h *ProductHandler) GetBrands(c *gin.Context) {
	brands, err := h.productService.GetBrands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

// @Summary Получить продукт по ID
// @Description Получить информацию о конкретном продукте
// @Tags products
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]interface{}
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// @Summary Создать продукт
// @Description Создать новый продукт (только для админов)
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Данные продукта"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.CreateProduct(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

// @Summary Обновить продукт
// @Description Обновить информацию о продукте (только для админов)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Param product body models.Product true "Обновленные данные продукта"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.UpdateProduct(uint(productID), &product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

// @Summary Удалить продукт
// @Description Удалить продукт (только для админов)
// @Tags products
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}