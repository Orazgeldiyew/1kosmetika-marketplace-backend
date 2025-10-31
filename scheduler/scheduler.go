package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"1kosmetika-marketplace-backend/repositories"
)


func StartCronJobs() {
	c := cron.New()


	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("📊 Running daily stats job...")

		statsRepo := repositories.NewStatsRepository()
		err := statsRepo.SaveDailyStats(time.Now())
		if err != nil {
			log.Println("❌ Failed to save daily stats:", err)
		} else {
			log.Println("✅ Daily stats saved successfully")
		}
	})
	if err != nil {
		log.Println("❌ Failed to schedule daily stats job:", err)
	}


	_, err = c.AddFunc("@hourly", func() {
		log.Println("♻️ Hourly cache refresh job running...")
	})
	if err != nil {
		log.Println("❌ Failed to schedule hourly job:", err)
	}

	c.Start()
	log.Println("🚀 Cron scheduler started")
}
