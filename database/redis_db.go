package database

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

var (
	redisInstance *Redis
	redisOnce     sync.Once
	redisTestMode bool
	redisProdMode bool
	redisDevMode  bool
	redisTestAddr string
)

func SetRedisTestMode(addr string) {
	redisTestMode = true
	redisProdMode = false
	redisDevMode = false
	redisTestAddr = addr
	redisInstance = nil
	redisOnce = sync.Once{}
}

func SetRedisDevMode() {
	redisTestMode = false
	redisProdMode = false
	redisDevMode = true
	redisTestAddr = ""
	redisInstance = nil
	redisOnce = sync.Once{}
}

func GetRedis() *Redis {
	redisOnce.Do(func() {
		redisInstance = &Redis{}
		if redisTestMode {
			redisInstance.initialiseTestRedis(redisTestAddr)
		}
		if redisDevMode {
			redisInstance.initialiseDevRedis()
		}
		if redisProdMode {
			redisInstance.initialiseProdRedis()
		}
		if !redisDevMode && !redisTestMode && !redisProdMode {
			log.Fatalf("No Redis Mode Set! :(")
		}
	})
	return redisInstance
}

func (r *Redis) initialiseTestRedis(addr string) {
	r.client = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	log.Printf("Redis Test Connection Succeeded")
}

func (r *Redis) initialiseDevRedis() {
	err := godotenv.Load("/home/callum/Desktop/Go-Web-App-Backend/.dev.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default fallback
	}

	r.client = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	log.Printf("Redis Dev Connection Succeeded")
}

func (r *Redis) initialiseProdRedis() {
	err := godotenv.Load("/home/callum/Desktop/Go-Web-App-Backend/.prod.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // default fallback
	}

	r.client = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	log.Printf("Redis Dev Connection Succeeded")
}

func (r *Redis) ClearAll() error {
	return r.client.FlushAll(context.Background()).Err()
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}
