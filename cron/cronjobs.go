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
		services.SchedulerEmail()
	})

	c.Start()

	// Keep the application running
	select {}
}
