package service

import (
	"context"
	"errors"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
)

type UserService interface {
	Create(ctx context.Context, email, username, password, firstName, lastName string) (*models.User, error)
	GetAll(ctx context.Context, params utils.PaginationParams) ([]models.UserResponse, utils.PaginationMeta, error)
	GetByID(ctx context.Context, id interface{}) (*models.UserResponse, error)
	Update(ctx context.Context, id interface{}, req models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id interface{}) error
}

type userService struct {
	userRepo repository.UserRepository
	cartRepo repository.CartRepository
}

func NewUserService(userRepo repository.UserRepository, cartRepo repository.CartRepository) UserService {
	return &userService{
		userRepo: userRepo,
		cartRepo: cartRepo,
	}
}

func (s *userService) Create(ctx context.Context, email, username, password, firstName, lastName string) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	// Create cart for user (using cartRepo later, for now we will assume it's there or just use GORM directly if not available yet)
	// For now, let's just use the user object
	return &user, nil
}

func (s *userService) GetAll(ctx context.Context, params utils.PaginationParams) ([]models.UserResponse, utils.PaginationMeta, error) {
	users, totalCount, err := s.userRepo.GetAll(ctx, params)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CalculatePaginationMeta(params, totalCount)
	return users, meta, nil
}

func (s *userService) GetByID(ctx context.Context, id interface{}) (*models.UserResponse, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id interface{}, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id interface{}) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, user)
}
