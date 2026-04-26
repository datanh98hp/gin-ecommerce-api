package routes

import (
	"github.com/dat19/gin-ecommerce-api/internal/api/handlers"
	"github.com/dat19/gin-ecommerce-api/internal/api/middleware"
	"github.com/dat19/gin-ecommerce-api/internal/config"
	"github.com/dat19/gin-ecommerce-api/internal/database"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	productRepo := repository.NewProductRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	cartRepo := repository.NewCartRepository(database.DB)
	postRepo := repository.NewPostRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)

	// Initialize services
	productSvc := service.NewProductService(productRepo)
	userSvc := service.NewUserService(userRepo, cartRepo)
	authSvc := service.NewAuthService(cfg, userRepo, cartRepo)
	postSvc := service.NewPostService(postRepo)
	cartSvc := service.NewCartService(cartRepo, productRepo)
	orderSvc := service.NewOrderService(orderRepo, cartRepo, productRepo)
	seedSvc := service.NewSeedService(database.DB, userRepo, productRepo, cartRepo, postRepo, orderRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authSvc, userSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	productHandler := handlers.NewProductHandler(productSvc)
	postHandler := handlers.NewPostHandler(postSvc)
	cartHandler := handlers.NewCartHandler(cartSvc)
	orderHandler := handlers.NewOrderHandler(orderSvc)
	seedHandler := handlers.NewSeedHandler(seedSvc)

	// Apply global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public product routes
		v1.GET("/products", productHandler.GetAll)
		v1.GET("/products/:id", productHandler.GetByID)

		// Public post routes
		v1.GET("/posts", postHandler.GetAll)
		v1.GET("/posts/:id", postHandler.GetByID)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Auth routes
			protected.GET("/auth/me", authHandler.Me)
			protected.POST("/auth/logout", authHandler.Logout)

			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetAll)
				users.GET("/:id", userHandler.GetByID)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", middleware.AdminMiddleware(), userHandler.Delete)
			}

			// Product management (admin only)
			products := protected.Group("/products")
			products.Use(middleware.AdminMiddleware())
			{
				products.POST("", productHandler.Create)
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.Delete)
			}

			// Post routes
			posts := protected.Group("/posts")
			{
				posts.POST("", postHandler.Create)
				posts.PUT("/:id", postHandler.Update)
				posts.DELETE("/:id", postHandler.Delete)
			}

			// Cart routes
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)
				cart.POST("/items", cartHandler.AddItem)
				cart.PUT("/items/:itemId", cartHandler.UpdateItem)
				cart.DELETE("/items/:itemId", cartHandler.RemoveItem)
				cart.DELETE("", cartHandler.ClearCart)
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.Create)
				orders.GET("", orderHandler.GetAll)
				orders.GET("/:id", orderHandler.GetByID)
				orders.POST("/:id/cancel", orderHandler.Cancel)
				orders.PUT("/:id/status", middleware.AdminMiddleware(), orderHandler.UpdateStatus)
			}

			// Seed routes (admin only)
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("/seed", seedHandler.SeedData)
			}
		}
	}
}
