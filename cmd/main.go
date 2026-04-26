package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dat19/gin-ecommerce-api/internal/api/routes"
	"github.com/dat19/gin-ecommerce-api/internal/cache"
	"github.com/dat19/gin-ecommerce-api/internal/config"
	"github.com/dat19/gin-ecommerce-api/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	// Connect to Redis
	if err := cache.Connect(cfg); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	// Initialize Gin engine
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, cfg)

	// Start server
	go func() {
		addr := ":" + cfg.Server.Port
		log.Printf("Server starting on %s in %s mode", addr, cfg.Server.Env)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	if err := cache.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}
	log.Println("Server stopped")
}
