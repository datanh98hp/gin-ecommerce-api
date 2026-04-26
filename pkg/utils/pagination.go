package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
}

// PaginationMeta holds pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
}

// GetPaginationParams extracts pagination parameters from query params
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// Validate page and page_size
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	// Validate order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

// Paginate applies pagination to a GORM query
func Paginate(params PaginationParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.PageSize
		orderBy := params.Sort + " " + params.Order
		return db.Offset(offset).Limit(params.PageSize).Order(orderBy)
	}
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(params PaginationParams, totalCount int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))
	return PaginationMeta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
