package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	CTX    = context.Background()
	Client *redis.Client
)

func Connect() error {
	Client = redis.NewClient(&redis.Options{Addr: "localhost:6381"})
	return Client.Ping(CTX).Err()
}

func Enqueue(jobID int64) error {
	return Client.RPush(CTX, "jobs", jobID).Err()
}

func Increment(key string) error {
	return Client.IncrBy(CTX, key, 1).Err()
}

func GetValue(key string) (string, error) {
	return Client.Get(CTX, key).Result()
}
func SetKey(key string) error {
	return Client.Set(CTX, key, 0, 10*time.Second).Err()
}

func CheckAllowed(key string) (bool, error) {
	count, err := Client.Incr(CTX, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := Client.Expire(CTX, key, 10*time.Second).Err(); err != nil {
			return false, err
		}
	}

	return count <= 5, nil
}
