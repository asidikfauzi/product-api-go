package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestInitRedis(t *testing.T) {
	mockRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start mock Redis server: %v", err)
	}
	defer mockRedis.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
		DB:   0,
	})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	assert.NoError(t, err, "Failed to connect to mock Redis")

	err = redisClient.Set(ctx, "test-key", "test-value", 0).Err()
	assert.NoError(t, err, "Failed to set value in Redis")

	val, err := redisClient.Get(ctx, "test-key").Result()
	assert.NoError(t, err, "Failed to get value from Redis")
	assert.Equal(t, "test-value", val, "Value mismatch")

	fmt.Println("Redis test passed!")
}
