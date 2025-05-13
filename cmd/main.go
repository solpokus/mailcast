package main

import (
	"mailcast/configuration"
	"mailcast/cron"
	"mailcast/database"
)

func main() {

	// Load configuration
	cfg := configuration.LoadConfig()

	// Initialize database connection
	database.ConnectGORM(cfg)

	// Add a 10-second delay
	// fmt.Println("Waiting for 10 seconds...")
	// time.Sleep(10 * time.Second)

	// checkEmailAndStart()

	cron.SchedEmail()

	// services.CheckEmailOauthAndStart()

	// repository.SyncTaskRedisToPostgre()
	// repository.InsertTaskToRedis()

	// client.MainClient()

	// tasks.MainAsynq()

	// parser()

	// user, err := services.GetUserByEmailService("admin@gmail.com")
	// Fetch User
	// var user models.User
	// result := database.DB.Debug().Unscoped().Where("email = ?", "test3@gmail.com").First(&user)

	// if result.Error != nil {
	// 	log.Println("❌ Error fetching user:", result.Error)
	// } else {
	// 	fmt.Println("✅ User found: ", user.FirstName)
	// }

}
