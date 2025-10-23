package repositories

import (
	"1kosmetika-marketplace-backend/database"
	"log"
	"time"
)

// CleanupOldLogs удаляет записи логов старше X дней
func CleanupOldLogs(days int) error {
	db := database.DB
	threshold := time.Now().AddDate(0, 0, -days)

	result := db.Exec("DELETE FROM logs WHERE created_at < ?", threshold)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("🗑️ Deleted %d old logs (older than %d days)\n", result.RowsAffected, days)
	return nil
}
