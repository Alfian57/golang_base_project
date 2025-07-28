package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	// Register the channel to receive specific signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine so it doesn't block
	go func() {
		logger.Log.Infof("Server starting on %s", cfg.Server.Url)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Server failed to start: %v", err)
		}
		logger.Log.Info("Server stopped listening")
	}()

	logger.Log.Info("Server is ready to handle requests. Press Ctrl+C to shutdown")

	// Wait for interrupt signal to gracefully shutdown the server
	sig := <-quit
	logger.Log.Infof("Received signal: %v. Shutting down server...", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Disable keep-alives to help with graceful shutdown
	server.SetKeepAlivesEnabled(false)

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Log.Info("Server gracefully stopped")
	os.Exit(0)
}
