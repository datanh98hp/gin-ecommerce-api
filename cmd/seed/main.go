package main

import (
	"context"
	"log"

	"github.com/dat19/gin-ecommerce-api/internal/config"
	"github.com/dat19/gin-ecommerce-api/internal/database"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	cartRepo := repository.NewCartRepository(database.DB)
	postRepo := repository.NewPostRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)

	// Initialize SeedService
	seedSvc := service.NewSeedService(database.DB, userRepo, productRepo, cartRepo, postRepo, orderRepo)

	// Run seeding
	ctx := context.Background()
	if err := seedSvc.SeedAll(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully")
}
