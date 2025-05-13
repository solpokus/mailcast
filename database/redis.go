package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GetAllAsynqTaskKeys scans Redis for keys matching "asynq:{default}:t:*"
func GetAllAsynqTaskKeys(ctx context.Context, rdb *redis.Client) ([]string, error) {
	var cursor uint64
	var keys []string

	for {
		k, newCursor, err := rdb.Scan(ctx, cursor, "asynq:{default}:t:*", 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
