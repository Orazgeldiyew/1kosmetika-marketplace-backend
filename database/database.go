package database

import (
	"fmt"
	"1kosmetika-marketplace-backend/config"
	"1kosmetika-marketplace-backend/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ashgabat",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
		&models.Cart{},
		&models.CartItem{},
		&models.Review{},
		&models.Favorite{},  
		&models.Notification{},
	)
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	log.Println("✅ Database migration completed")
	return nil
}