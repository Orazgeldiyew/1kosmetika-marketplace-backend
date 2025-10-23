package database

import (
	"fmt"
	"log"

	"1kosmetika-marketplace-backend/config"
	"1kosmetika-marketplace-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ashgabat",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("✅ Database connection established")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderProduct{}, // ✅ make sure join table with qty/price exists
		&models.Cart{},
		&models.CartItem{},
		&models.Review{},
		&models.Favorite{},
		&models.Notification{},
		&models.DailyStats{},
	)
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	log.Println("✅ Database migration completed")
	return nil
}
	