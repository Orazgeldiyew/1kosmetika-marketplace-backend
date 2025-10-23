package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadProductImageHandler — загрузка фото товара
func UploadProductImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл"})
		return
	}

	ct := file.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Разрешены только изображения"})
		return
	}

	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании директории"})
		return
	}

	filename := uuid.NewString() + filepath.Ext(file.Filename)
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
		return
	}

	fileURL := "/static/products/" + filename

	c.JSON(http.StatusOK, gin.H{
		"message": "Файл успешно загружен ✅",
		"url":     fileURL,
		"ts":      time.Now().Unix(),
	})
}
