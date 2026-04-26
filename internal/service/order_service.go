package service

import (
	"context"
	"errors"

	"github.com/dat19/gin-ecommerce-api/internal/database"
	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type OrderService interface {
	Create(ctx context.Context, userID uint, req models.CreateOrderRequest) (*models.Order, error)
	GetAll(ctx context.Context, userID uint, role string, params utils.PaginationParams) ([]models.Order, utils.PaginationMeta, error)
	GetByID(ctx context.Context, id interface{}, userID uint, role string) (*models.Order, error)
	Cancel(ctx context.Context, id interface{}, userID uint) error
	UpdateStatus(ctx context.Context, id interface{}, status string) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) Create(ctx context.Context, userID uint, req models.CreateOrderRequest) (*models.Order, error) {
	// Get cart
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil || len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var order *models.Order

	// Transactional create
	err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []models.OrderItem

		for _, cartItem := range cart.CartItems {
			// Check stock again
			product := cartItem.Product
			if product.Stock < cartItem.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}

			// Deduct stock
			product.Stock -= cartItem.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			subtotal := product.Price * float64(cartItem.Quantity)
			totalAmount += subtotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  cartItem.Quantity,
				Price:     product.Price,
			})
		}

		// Create order
		order = &models.Order{
			UserID:        userID,
			TotalAmount:   totalAmount,
			ShippingAddr:  req.ShippingAddress,
			PaymentMethod: req.PaymentMethod,
			Status:        "pending",
			OrderItems:    orderItems,
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Clear cart
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(ctx, order.ID)
}

func (s *orderService) GetAll(ctx context.Context, userID uint, role string, params utils.PaginationParams) ([]models.Order, utils.PaginationMeta, error) {
	isAdmin := role == "admin"
	orders, totalCount, err := s.orderRepo.GetAll(ctx, params, userID, isAdmin)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CalculatePaginationMeta(params, totalCount)
	return orders, meta, nil
}

func (s *orderService) GetByID(ctx context.Context, id interface{}, userID uint, role string) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.ErrNotFound
	}

	if role != "admin" && order.UserID != userID {
		return nil, utils.ErrForbidden
	}

	return order, nil
}

func (s *orderService) Cancel(ctx context.Context, id interface{}, userID uint) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return utils.ErrNotFound
	}

	if order.UserID != userID {
		return utils.ErrForbidden
	}

	if order.Status != "pending" {
		return errors.New("can only cancel pending orders")
	}

	order.Status = "cancelled"
	return s.orderRepo.Update(ctx, order)
}

func (s *orderService) UpdateStatus(ctx context.Context, id interface{}, status string) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return utils.ErrNotFound
	}

	order.Status = status
	return s.orderRepo.Update(ctx, order)
}
