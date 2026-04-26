package repository

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetAll(ctx context.Context, params utils.PaginationParams, category, productType string) ([]models.Product, int64, error)
	GetByID(ctx context.Context, id interface{}) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetAll(ctx context.Context, params utils.PaginationParams, category, productType string) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if productType != "" {
		query = query.Where("type = ?", productType)
	}

	query = query.Where("is_active = ?", true)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(utils.Paginate(params)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func (r *productRepository) GetByID(ctx context.Context, id interface{}) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Delete(product).Error
}
