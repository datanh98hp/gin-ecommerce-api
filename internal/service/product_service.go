package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dat19/gin-ecommerce-api/internal/cache"
	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/repository"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
)

type ProductService interface {
	Create(ctx context.Context, req models.CreateProductRequest) (*models.Product, error)
	GetAll(ctx context.Context, params utils.PaginationParams, category, productType string) ([]models.Product, utils.PaginationMeta, error)
	GetByID(ctx context.Context, id interface{}) (*models.Product, error)
	Update(ctx context.Context, id interface{}, req models.UpdateProductRequest) (*models.Product, error)
	Delete(ctx context.Context, id interface{}) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Type:        req.Type,
		ImageURL:    req.ImageURL,
	}

	if err := s.repo.Create(ctx, &product); err != nil {
		return nil, err
	}

	// Invalidate cache
	cache.DeleteByPrefix(ctx, "products:list")

	return &product, nil
}

func (s *productService) GetAll(ctx context.Context, params utils.PaginationParams, category, productType string) ([]models.Product, utils.PaginationMeta, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf("products:list:page=%d:size=%d:cat=%s:type=%s:sort=%s:order=%s",
		params.Page, params.PageSize, category, productType, params.Sort, params.Order)

	// Try to get from cache
	type CachedResult struct {
		Products []models.Product     `json:"products"`
		Meta     utils.PaginationMeta `json:"meta"`
	}

	var cached CachedResult
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached.Products, cached.Meta, nil
	}

	// Fetch from repository
	products, totalCount, err := s.repo.GetAll(ctx, params, category, productType)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.CalculatePaginationMeta(params, totalCount)

	// Store in cache
	cacheData := CachedResult{
		Products: products,
		Meta:     meta,
	}
	_ = cache.Set(ctx, cacheKey, cacheData, 5*time.Minute)

	return products, meta, nil
}

func (s *productService) GetByID(ctx context.Context, id interface{}) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) Update(ctx context.Context, id interface{}, req models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Type != nil {
		product.Type = *req.Type
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Invalidate cache
	cache.DeleteByPrefix(ctx, "products:list")

	return product, nil
}

func (s *productService) Delete(ctx context.Context, id interface{}) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, product); err != nil {
		return err
	}

	// Invalidate cache
	cache.DeleteByPrefix(ctx, "products:list")

	return nil
}
