package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
}

func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func SetupDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "pgadmin"),
		GetEnv("DB_NAME", "wallet-topup"),
		GetEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	return db
}

func SetupRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: GetEnv("REDIS_ADDR", "localhost:6379"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("❌ Failed to connect Redis:", err)
	}
	return client
}
