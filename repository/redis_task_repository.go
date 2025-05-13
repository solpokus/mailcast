package repository

import (
	"context"
	"fmt"
	"log"
	"mailcast/configuration"
	"mailcast/database"
	"mailcast/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Synchronize Task Redis to Postgre
func SyncTaskRedisToPostgre() {
	log.Println("Start SyncTaskRedisToPostgre")
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: configuration.CONFIG.RedisAddr, // ✅ your Redis address
	})

	keys, err := database.GetAllAsynqTaskKeys(ctx, rdb)
	if err != nil {
		log.Println("Error fetching task keys:", err)
		return
	}

	log.Println("Length : %v", len(keys))

	for _, key := range keys {
		log.Println("=== Start Found task key:", key)

		// Insert id to postgre

		log.Println("=== Start insert to postgre key:", key)
		IdNewRedisTask := InsertNewRedisTask(key)
		log.Println("=== Finish insert to postgre key:", key)

		// HGETALL returns all fields in the hash stored at key
		result, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			panic(err)
		}

		log.Println("--- Task Data :")
		for field, value := range result {
			// log.Printf("%s: %s\n", field, value)
			if IdNewRedisTask != 0 {

				var jsonStr string
				if strings.Contains(value, "type:notif") {

					jsonStr = extractJSON(value)

					if jsonStr == "" {
						log.Fatal("❌ Failed to extract JSON from string")
					}

					// fmt.Println("✅ Extracted JSON:")
					// fmt.Println(jsonStr)
					// log.Printf("%s: %s\n", field, value)
					InsertNewRedisTaskDetail(IdNewRedisTask, field, jsonStr)
				} else {

					log.Printf("%s: %s\n", field, value)
					InsertNewRedisTaskDetail(IdNewRedisTask, field, value)
				}
			}
		}

		log.Println("==== End Found task key:", key)
	}

	log.Println("Finish SyncTaskRedisToPostgre")
}

// Inserts a new redis_task into the postgres database
func InsertNewRedisTask(idAsynq string) uint {

	// Insert a new record
	newMessage := models.RedisTask{
		IdAsynq: idAsynq,
	}

	// Save to database
	// result := database.DB.Debug().Create(&newMessage)
	result := database.DB.Create(&newMessage)
	if result.Error != nil {
		log.Fatalf("Error inserting redis_task : %v", result.Error)
	}

	// Print inserted record ID
	log.Println("✅ Message inserted to redis_task with ID:", newMessage.ID)

	return newMessage.ID
}

// extractJSON tries to extract the first JSON object from a string.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start : end+1]
}

func InsertTaskToRedis() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: configuration.CONFIG.RedisAddr, // ✅ your Redis address
	})

	key := "asynq:{default}:t:tttt"

	fields := map[string]interface{}{
		"msg":   `{"Payload":{"messages":[{"caption":"Dear RAMA, JUAN MR, \n\nThis is a reminder..."}]}}`,
		"state": "completed",
	}

	err := rdb.HSet(ctx, key, fields).Err()
	if err != nil {
		panic(fmt.Sprintf("❌ Redis HSET failed: %v", err))
	}

	fmt.Println("✅ Task inserted with multiple fields.")
}
