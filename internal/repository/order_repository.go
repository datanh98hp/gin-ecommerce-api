package repository

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetAll(ctx context.Context, params utils.PaginationParams, userID uint, isAdmin bool) ([]models.Order, int64, error)
	GetByID(ctx context.Context, id interface{}) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetAll(ctx context.Context, params utils.PaginationParams, userID uint, isAdmin bool) ([]models.Order, int64, error) {
	var orders []models.Order
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Preload("OrderItems.Product").Preload("User")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(utils.Paginate(params)).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id interface{}) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Product").Preload("User").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
