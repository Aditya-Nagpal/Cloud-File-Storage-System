package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/auth-service/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisURL,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatal("Error connecting to redis database: ", err)
	}

	fmt.Println("Connected to redis successfully")
}
