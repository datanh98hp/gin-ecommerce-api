package repository

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetAll(ctx context.Context, params utils.PaginationParams) ([]models.UserResponse, int64, error)
	GetByID(ctx context.Context, id interface{}) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id interface{}) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetAll(ctx context.Context, params utils.PaginationParams) ([]models.UserResponse, int64, error) {
	var users []models.User
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(utils.Paginate(params)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	list := make([]models.UserResponse, len(users))
	for i, user := range users {
		list[i] = mappingUserResponse(&user)
	}
	return list, totalCount, nil
}
func mappingUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id interface{}) (*models.UserResponse, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	res := mappingUserResponse(&user)
	return &res, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.UserResponse, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	res := models.UserResponse{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}
	return &res, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id interface{}) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
