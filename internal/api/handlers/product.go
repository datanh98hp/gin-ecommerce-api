package handlers

import (
	"net/http"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/dat19/gin-ecommerce-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateCreateProductRequest(req.Name, req.Description, req.Price, req.Stock, req.Category, req.Type, req.ImageURL); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Normalize data
	req.Name = validator.TrimWhitespace(req.Name)
	req.Description = validator.TrimWhitespace(req.Description)

	product, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	params := utils.GetPaginationParams(c)
	category := c.Query("category")
	productType := c.Query("type")

	products, meta, err := h.service.GetAll(c.Request.Context(), params, category, productType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, "Products retrieved successfully", products, meta)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateUpdateProductRequest(req.Name, req.Description, req.Price, req.Stock, req.Category, req.Type, req.ImageURL, req.IsActive); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}
