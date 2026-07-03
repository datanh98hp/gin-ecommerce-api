package handlers

import (
	"net/http"
	"strconv"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/dat19/gin-ecommerce-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(svc service.CartService) *CartHandler {
	return &CartHandler{service: svc}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	cart, err := h.service.GetCart(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve cart")
		return
	}

	// Build response
	var cartItems []models.CartItemResponse
	var totalPrice float64

	for _, item := range cart.CartItems {
		subtotal := item.Product.Price * float64(item.Quantity)
		totalPrice += subtotal

		cartItems = append(cartItems, models.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	response := models.CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		CartItems:  cartItems,
		TotalPrice: totalPrice,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", response)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateAddToCartRequest(req.ProductID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AddItem(c.Request.Context(), userID, req.ProductID, req.Quantity); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrNotFound {
			status = http.StatusNotFound
		} else if err.Error() == "insufficient stock" || err.Error() == "insufficient stock for total quantity" {
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item added to cart successfully", nil)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	itemID := parseUint(c.Param("itemId"))
	val, _ := c.Get("user_id")
	userID := val.(uint)

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateUpdateCartItemRequest(req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate cart item ID
	if err := validator.ValidateCartItemID(itemID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateItem(c.Request.Context(), userID, itemID, req.Quantity); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrForbidden {
			status = http.StatusForbidden
		} else if err == utils.ErrNotFound {
			status = http.StatusNotFound
		} else if err.Error() == "insufficient stock" {
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart item updated successfully", nil)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	itemID := parseUint(c.Param("itemId"))
	val, _ := c.Get("user_id")
	userID := val.(uint)

	// Validate cart item ID
	if err := validator.ValidateCartItemID(itemID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.RemoveItem(c.Request.Context(), userID, itemID); err != nil {
		status := http.StatusInternalServerError
		if err == utils.ErrForbidden {
			status = http.StatusForbidden
		} else if err == utils.ErrNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item removed from cart successfully", nil)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	if err := h.service.ClearCart(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear cart")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}

func parseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 32)
	return uint(val)
}
