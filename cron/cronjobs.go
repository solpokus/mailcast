package cron

import (
	"fmt"
	"mailcast/services"
	"time"

	"github.com/robfig/cron"
)

func SchedEmail() {
	c := cron.New()

	// Schedule a task every 1 minutes
	c.AddFunc("0 */1 * * * *", func() {
		fmt.Println("⏳ Running scheduled task at:", time.Now())

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("⚠️ Panic recovered in scheduled task: %v\n", r)
			}
		}()

		if err := services.SchedulerEmail(); err != nil {
			fmt.Printf("❌ Error in SchedulerEmail: %v\n", err)
			// Optionally: log to file or external monitor here
		}
	})

	c.Start()

	// Keep the application running
	select {}
}
