package repository

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	Create(ctx context.Context, cart *models.Cart) error
	GetByUserID(ctx context.Context, userID uint) (*models.Cart, error)
	GetByID(ctx context.Context, id uint) (*models.Cart, error)
	Update(ctx context.Context, cart *models.Cart) error
	Clear(ctx context.Context, cartID uint) error
	// Item methods
	GetItemByID(ctx context.Context, itemID uint) (*models.CartItem, error)
	AddItem(ctx context.Context, item *models.CartItem) error
	UpdateItem(ctx context.Context, item *models.CartItem) error
	RemoveItem(ctx context.Context, itemID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByID(ctx context.Context, id uint) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).Preload("Items.Product").First(&cart, id).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *cartRepository) Clear(ctx context.Context, cartID uint) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetItemByID(ctx context.Context, itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	if err := r.db.WithContext(ctx).First(&item, itemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) AddItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *cartRepository) UpdateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, itemID uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, itemID).Error
}
