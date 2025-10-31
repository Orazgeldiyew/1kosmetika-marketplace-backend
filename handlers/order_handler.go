package handlers

import (
	"1kosmetika-marketplace-backend/services"
	"1kosmetika-marketplace-backend/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderRequest struct {
	ProductIDs []uint `json:"product_ids" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data: " + err.Error()})
		return
	}

	if len(req.ProductIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product list cannot be empty"})
		return
	}

	userID := c.GetUint("user_id")
	order, err := h.orderService.CreateOrder(userID, req.ProductIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userEmail := c.GetString("user_email") // если ты сохраняешь email в токене
	if userEmail == "" {
		userEmail = "test@example.com" // временно, чтобы проверить работу
	}

	subject := "Ваш заказ успешно создан!"
	body := fmt.Sprintf(`
	<h2>Здравствуйте!</h2>
	<p>Ваш заказ #%d успешно оформлен и ожидает подтверждения.</p>
	<p>Общая сумма заказа: <b>%.2f</b></p>
`, order.ID, order.Total)

	if err := utils.SendEmail(userEmail, subject, body); err != nil {
		fmt.Println("❌ Ошибка при отправке email:", err)
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
