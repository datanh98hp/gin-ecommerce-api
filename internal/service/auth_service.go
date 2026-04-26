package service

import (
	"context"
	"errors"

	"github.com/dat19/gin-ecommerce-api/internal/config"
	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
}

type authService struct {
	cfg      *config.Config
	userRepo repository.UserRepository
	cartRepo repository.CartRepository
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, cartRepo repository.CartRepository) AuthService {
	return &authService{
		cfg:      cfg,
		userRepo: userRepo,
		cartRepo: cartRepo,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	_, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	// Create cart for user
	cart := models.Cart{UserID: user.ID}
	_ = s.cartRepo.Create(ctx, &cart)

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireTime)
	if err != nil {
		return nil, err
	}
	
	return &models.AuthResponse{
		Token: token,
		User:  models.UserResponse{
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
	},
	}, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireTime)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  models.UserResponse{
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
		},
	}, nil
}
