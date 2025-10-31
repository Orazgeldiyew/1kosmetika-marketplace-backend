package database

import (
	"fmt"
	"log"

	"1kosmetika-marketplace-backend/config"
	"1kosmetika-marketplace-backend/models"

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

	DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='order_products' AND column_name='price'
			) THEN
				ALTER TABLE order_products ADD COLUMN price numeric DEFAULT 0;
			END IF;
		END$$;
	`)


	DB.Exec(`UPDATE order_products SET price = 0 WHERE price IS NULL;`)
	DB.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='order_products' AND column_name='price'
			) THEN
				ALTER TABLE order_products ALTER COLUMN price SET NOT NULL;
			END IF;
		END$$;
	`)


	DB.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='order_products' AND column_name='quantity'
			) THEN
				UPDATE order_products SET quantity = 1 WHERE quantity IS NULL;
				ALTER TABLE order_products ALTER COLUMN quantity SET DEFAULT 1;
				ALTER TABLE order_products ALTER COLUMN quantity SET NOT NULL;
			END IF;
		END$$;
	`)


	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderProduct{},
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
