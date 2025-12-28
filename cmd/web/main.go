package main

import (
	"TENDABOX/routes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			getEnv("REDIS_HOST", "redis"),
			getEnv("REDIS_PORT", "6379"),
		),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("redis", rdb)
		c.Next()
	})

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// ⭐ all routes here
	routes.RegisterRoutes(r, db)

	fmt.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
