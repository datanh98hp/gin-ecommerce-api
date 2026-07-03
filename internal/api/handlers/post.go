package handlers

import (
	"net/http"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/dat19/gin-ecommerce-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{service: svc}
}

func (h *PostHandler) Create(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateCreatePostRequest(req.Title, req.Content); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Normalize data
	req.Title = validator.TrimWhitespace(req.Title)
	req.Content = validator.TrimWhitespace(req.Content)

	post, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Post created successfully", post)
}

func (h *PostHandler) GetAll(c *gin.Context) {
	params := utils.GetPaginationParams(c)

	posts, meta, err := h.service.GetAll(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, "Posts retrieved successfully", posts, meta)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	post, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Post not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post retrieved successfully", post)
}

func (h *PostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	val, _ := c.Get("user_id")
	userID := val.(uint)
	role, _ := c.Get("user_role")

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateUpdatePostRequest(req.Title, req.Content, req.IsActive); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.service.Update(c.Request.Context(), id, userID, role.(string), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrForbidden {
			status = http.StatusForbidden
		} else if err == utils.ErrNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post updated successfully", post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	val, _ := c.Get("user_id")
	userID := val.(uint)
	role, _ := c.Get("user_role")

	if err := h.service.Delete(c.Request.Context(), id, userID, role.(string)); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrForbidden {
			status = http.StatusForbidden
		} else if err == utils.ErrNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Post deleted successfully", nil)
}
