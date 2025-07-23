package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Alfian57/belajar-golang/internal/config"
	"github.com/Alfian57/belajar-golang/internal/database"
	"github.com/Alfian57/belajar-golang/internal/logger"
	"github.com/Alfian57/belajar-golang/internal/middleware"
	"github.com/Alfian57/belajar-golang/internal/router"
	"github.com/Alfian57/belajar-golang/internal/validation"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger.Init()
	database.Init(cfg.Database)
	validation.Init()

	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}
	f, err := os.Create("logs/gin.log")
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}
	gin.DefaultWriter = io.MultiWriter(f)
	gin.DisableConsoleColor()

	router := router.NewRouter()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorMiddleware())

	if err := router.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
		panic(fmt.Sprintf("Failed to set trusted proxies: %v", err))
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowCredentials: cfg.Cors.AllowCredentials,
	}))

	server := &http.Server{
		Addr:    cfg.Server.Url,
		Handler: router,
	}

	logger.Log.Infof("Server starting on %s", cfg.Server.Url)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatalf("Server failed to start: %v", err)
	}
}
