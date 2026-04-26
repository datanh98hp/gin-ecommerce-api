package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SeedService interface {
	SeedAll(ctx context.Context) error
	SeedUsers(ctx context.Context, count int) ([]models.User, error)
	SeedProducts(ctx context.Context, count int) ([]models.Product, error)
	SeedPosts(ctx context.Context, count int, users []models.User) error
	SeedOrders(ctx context.Context, count int, users []models.User, products []models.Product) error
}

type seedService struct {
	db          *gorm.DB
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	cartRepo    repository.CartRepository
	postRepo    repository.PostRepository
	orderRepo   repository.OrderRepository
}

func NewSeedService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
	cartRepo repository.CartRepository,
	postRepo repository.PostRepository,
	orderRepo repository.OrderRepository,
) SeedService {
	return &seedService{
		db:          db,
		userRepo:    userRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		postRepo:    postRepo,
		orderRepo:   orderRepo,
	}
}

func (s *seedService) SeedAll(ctx context.Context) error {
	users, err := s.SeedUsers(ctx, 10)
	if err != nil {
		return err
	}

	products, err := s.SeedProducts(ctx, 50)
	if err != nil {
		return err
	}

	if err := s.SeedPosts(ctx, 30, users); err != nil {
		return err
	}

	if err := s.SeedOrders(ctx, 15, users, products); err != nil {
		return err
	}

	return nil
}

func (s *seedService) SeedUsers(ctx context.Context, count int) ([]models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := make([]models.User, 0, count)

	for i := 0; i < count; i++ {
		username := fmt.Sprintf("user_%d_%d", i+1, rand.Intn(1000))
		user := models.User{
			Username:  username,
			Email:     fmt.Sprintf("%s@example.com", username),
			Password:  string(hashedPassword),
			FirstName: fmt.Sprintf("First%d", i+1),
			LastName:  fmt.Sprintf("Last%d", i+1),
			Role:      "user",
			IsActive:  true,
		}

		if err := s.userRepo.Create(ctx, &user); err != nil {
			return nil, err
		}

		// Create cart
		cart := models.Cart{UserID: user.ID}
		s.cartRepo.Create(ctx, &cart)

		users = append(users, user)
	}

	return users, nil
}

func (s *seedService) SeedProducts(ctx context.Context, count int) ([]models.Product, error) {
	categories := []string{"Electronics", "Clothing", "Home & Garden", "Books", "Sports"}
	types := []string{"Premium", "Regular", "Budget"}
	products := make([]models.Product, 0, count)

	for i := 0; i < count; i++ {
		cat := categories[rand.Intn(len(categories))]
		product := models.Product{
			Name:        fmt.Sprintf("%s Product %d", cat, rand.Intn(10000)),
			Description: fmt.Sprintf("This is a high-quality product from our %s collection.", cat),
			Price:       float64(rand.Intn(100000)) / 100.0,
			Stock:       rand.Intn(100) + 10,
			Category:    cat,
			Type:        types[rand.Intn(len(types))],
			ImageURL:    fmt.Sprintf("https://picsum.photos/400/300?random=%d", rand.Intn(1000)),
			IsActive:    true,
		}

		if err := s.productRepo.Create(ctx, &product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *seedService) SeedPosts(ctx context.Context, count int, users []models.User) error {
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		post := models.Post{
			Title:    fmt.Sprintf("Announcement #%d", rand.Intn(1000)),
			Content:  "Detailed content for the announcement. Supporting Markdown and structured text.",
			UserID:   user.ID,
			IsActive: true,
		}

		if err := s.postRepo.Create(ctx, &post); err != nil {
			return err
		}
	}
	return nil
}

func (s *seedService) SeedOrders(ctx context.Context, count int, users []models.User, products []models.Product) error {
	statuses := []string{"pending", "processing", "shipped", "delivered"}
	address := "123 Main St, City, Country"

	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]

		// Start transaction via DB if needed, but here we use repos.
		// For seeding simplified, we just create.

		var totalAmount float64
		itemCount := rand.Intn(5) + 1
		orderItems := make([]models.OrderItem, 0, itemCount)

		for j := 0; j < itemCount; j++ {
			product := products[rand.Intn(len(products))]
			qty := rand.Intn(3) + 1
			orderItems = append(orderItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  qty,
				Price:     product.Price,
			})
			totalAmount += product.Price * float64(qty)
		}

		order := models.Order{
			UserID:        user.ID,
			Status:        statuses[rand.Intn(len(statuses))],
			TotalAmount:   totalAmount,
			ShippingAddr:  address,
			PaymentMethod: "Credit Card",
			OrderItems:    orderItems,
		}

		if err := s.orderRepo.Create(ctx, &order); err != nil {
			return err
		}
	}
	return nil
}
