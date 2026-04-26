package repository

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetAll(ctx context.Context, params utils.PaginationParams) ([]models.Post, int64, error)
	GetByID(ctx context.Context, id interface{}) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, post *models.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) GetAll(ctx context.Context, params utils.PaginationParams) ([]models.Post, int64, error) {
	var posts []models.Post
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Post{}).Preload("User").Where("is_active = ?", true)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(utils.Paginate(params)).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, totalCount, nil
}

func (r *postRepository) GetByID(ctx context.Context, id interface{}) (*models.Post, error) {
	var post models.Post
	if err := r.db.WithContext(ctx).Preload("User").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepository) Delete(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Delete(post).Error
}
