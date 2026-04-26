package handlers

import (
	"net/http"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	order, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Order created successfully", order)
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)
	roleVal, _ := c.Get("user_role")
	role := roleVal.(string)

	params := utils.GetPaginationParams(c)

	orders, meta, err := h.service.GetAll(c.Request.Context(), userID, role, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, "Orders retrieved successfully", orders, meta)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	val, _ := c.Get("user_id")
	userID := val.(uint)
	roleVal, _ := c.Get("user_role")
	role := roleVal.(string)

	order, err := h.service.GetByID(c.Request.Context(), id, userID, role)
	if err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrNotFound {
			status = http.StatusNotFound
		} else if err == utils.ErrForbidden {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateOrderStatusRequest
	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order status updated successfully", nil)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	val, _ := c.Get("user_id")
	userID := val.(uint)

	if err := h.service.Cancel(c.Request.Context(), id, userID); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrNotFound {
			status = http.StatusNotFound
		} else if err == utils.ErrForbidden {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order cancelled successfully", nil)
}
