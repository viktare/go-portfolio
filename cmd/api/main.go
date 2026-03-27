package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/viktare/go-portfolio/internal/config"
	"github.com/viktare/go-portfolio/internal/database"
	"github.com/viktare/go-portfolio/internal/router"
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

    r := gin.Default()
    r.SetTrustedProxies(nil)

    router.Setup(r, pool)

    r.Run(":" + cfg.Port)
}