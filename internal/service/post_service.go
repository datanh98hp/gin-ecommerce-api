package service

import (
	"context"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
)

type PostService interface {
	Create(ctx context.Context, userID uint, req models.CreatePostRequest) (*models.Post, error)
	GetAll(ctx context.Context, params utils.PaginationParams) ([]models.Post, utils.PaginationMeta, error)
	GetByID(ctx context.Context, id interface{}) (*models.Post, error)
	Update(ctx context.Context, id interface{}, userID uint, role string, req models.UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, id interface{}, userID uint, role string) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, userID uint, req models.CreatePostRequest) (*models.Post, error) {
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(ctx, &post); err != nil {
		return nil, err
	}

	// Preload user for response
	return s.repo.GetByID(ctx, post.ID)
}

func (s *postService) GetAll(ctx context.Context, params utils.PaginationParams) ([]models.Post, utils.PaginationMeta, error) {
	posts, totalCount, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CalculatePaginationMeta(params, totalCount)
	return posts, meta, nil
}

func (s *postService) GetByID(ctx context.Context, id interface{}) (*models.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *postService) Update(ctx context.Context, id interface{}, userID uint, role string, req models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Authorization: only author or admin can update
	if role != "admin" && post.UserID != userID {
		return nil, utils.ErrForbidden
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.IsActive != nil {
		post.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) Delete(ctx context.Context, id interface{}, userID uint, role string) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Authorization: only author or admin can delete
	if role != "admin" && post.UserID != userID {
		return utils.ErrForbidden
	}

	return s.repo.Delete(ctx, post)
}
