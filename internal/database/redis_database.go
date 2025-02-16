package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"product-api-go/internal/config"
	"strconv"
)

func InitRedis() *redis.Client {
	var ctx = context.Background()

	dbConfig := config.LoadRedisConfigFromEnv()
	dbName, _ := strconv.Atoi(dbConfig.DBName)
	dbPass := dbConfig.Password
	dbHost := dbConfig.Host
	dbPort := dbConfig.Port

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPass,
		DB:       dbName,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Connected to Redis!")
	return redisClient
}
