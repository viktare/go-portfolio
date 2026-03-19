package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/viktare/go-portfolio/internal/config"
	"github.com/viktare/go-portfolio/internal/database"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Server running",
		})
	})


	router.Run(":" + cfg.Port)
}