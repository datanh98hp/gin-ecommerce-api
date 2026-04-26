package service

import (
	"context"
	"errors"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
)

type CartService interface {
	GetCart(ctx context.Context, userID uint) (*models.Cart, error)
	AddItem(ctx context.Context, userID uint, productID uint, quantity int) error
	UpdateItem(ctx context.Context, userID uint, itemID uint, quantity int) error
	RemoveItem(ctx context.Context, userID uint, itemID uint) error
	ClearCart(ctx context.Context, userID uint) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If cart doesn't exist, create one
		cart = &models.Cart{UserID: userID}
		if err := s.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}
	return cart, nil
}

func (s *cartService) AddItem(ctx context.Context, userID uint, productID uint, quantity int) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return utils.ErrNotFound
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Check if item already in cart
	for _, item := range cart.CartItems {
		if item.ProductID == productID {
			item.Quantity += quantity
			if product.Stock < item.Quantity {
				return errors.New("insufficient stock for total quantity")
			}
			return s.cartRepo.UpdateItem(ctx, &item)
		}
	}

	// Add new item
	item := models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.AddItem(ctx, &item)
}

func (s *cartService) UpdateItem(ctx context.Context, userID uint, itemID uint, quantity int) error {
	item, err := s.cartRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return utils.ErrNotFound
	}

	cart, err := s.cartRepo.GetByID(ctx, item.CartID)
	if err != nil || cart.UserID != userID {
		return utils.ErrForbidden
	}

	product, err := s.productRepo.GetByID(ctx, item.ProductID)
	if err != nil {
		return utils.ErrNotFound
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	item.Quantity = quantity
	return s.cartRepo.UpdateItem(ctx, item)
}

func (s *cartService) RemoveItem(ctx context.Context, userID uint, itemID uint) error {
	item, err := s.cartRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return utils.ErrNotFound
	}

	cart, err := s.cartRepo.GetByID(ctx, item.CartID)
	if err != nil || cart.UserID != userID {
		return utils.ErrForbidden
	}

	return s.cartRepo.RemoveItem(ctx, itemID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uint) error {
	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return utils.ErrNotFound
	}
	return s.cartRepo.Clear(ctx, cart.ID)
}
