package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	CTX    = context.Background()
	Client *redis.Client
)

func Connect() error {
	Client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	return Client.Ping(CTX).Err()
}

func Enqueue(jobID int64) error {
	return Client.RPush(CTX, "jobs", jobID).Err()
}
