package validator

// USAGE_GUIDE.md - Practical Integration Examples for API Handlers

/*

# Validator Package Usage Guide

This guide shows how to integrate validators into your actual API handlers.

## 1. AUTH HANDLER EXAMPLE

```go
import (
	"gin-ecommerce-api/internal/models"
	"gin-ecommerce-api/pkg/validator"
)

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	
	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate using validator package
	if err := validator.ValidateRegisterRequest(
		req.Email, req.Username, req.Password, 
		req.FirstName, req.LastName,
	); err != nil {
		return c.JSON(400, gin.H{
			"error": err.Error(),
			"status": "validation_error",
		})
	}
	
	// Normalize data
	req.Email = validator.NormalizeEmail(req.Email)
	req.Username = validator.NormalizeUsername(req.Username)
	
	// Call service
	token, err := h.authService.Register(req)
	if err != nil {
		return c.JSON(409, gin.H{"error": "Email or username already exists"})
	}
	
	return c.JSON(201, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateLoginRequest(req.Email, req.Password); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Normalize
	req.Email = validator.NormalizeEmail(req.Email)
	
	// Process
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(401, gin.H{"error": "Invalid credentials"})
	}
	
	return c.JSON(200, gin.H{"token": token})
}
```

## 2. PRODUCT HANDLER EXAMPLE

```go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate using validator package
	if err := validator.ValidateCreateProductRequest(
		req.Name, req.Description, req.Price, req.Stock,
		req.Category, req.Type, req.ImageURL,
	); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Trim and normalize
	req.Name = validator.TrimWhitespace(req.Name)
	req.Description = validator.TrimWhitespace(req.Description)
	
	// Create product
	product, err := h.productService.Create(&req)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to create product"})
	}
	
	return c.JSON(201, gin.H{"data": product})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Invalid product ID"})
	}
	
	// Validate ID
	if err := validator.ValidateProductID(uint(id)); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	var req models.UpdateProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate using validator package
	if err := validator.ValidateUpdateProductRequest(
		req.Name, req.Description, req.Price, req.Stock,
		req.Category, req.Type, req.ImageURL, req.IsActive,
	); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Update product
	product, err := h.productService.Update(uint(id), &req)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to update product"})
	}
	
	return c.JSON(200, gin.H{"data": product})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	
	// Parse pagination params
	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)
	
	// Define allowed sort fields
	allowedFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name": true,
		"price": true,
		"stock": true,
	}
	
	// Validate and normalize
	params, err := validator.ValidateAndNormalizePaginationParams(
		p, ps, sortBy, sortOrder, allowedFields,
	)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Calculate offset
	offset := validator.CalculateOffset(params.Page, params.PageSize)
	limit := validator.GetLimit(params.PageSize)
	
	// Get products
	products, total, err := h.productService.List(offset, limit, params.SortBy, params.SortOrder)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to fetch products"})
	}
	
	return c.JSON(200, gin.H{
		"data": products,
		"pagination": gin.H{
			"page": params.Page,
			"page_size": params.PageSize,
			"total": total,
		},
	})
}
```

## 3. CART HANDLER EXAMPLE

```go
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req models.AddToCartRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateAddToCartRequest(req.ProductID, req.Quantity); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Get user ID from context
	userID, _ := c.Get("user_id")
	
	// Add to cart
	cartItem, err := h.cartService.AddToCart(userID.(uint), req.ProductID, req.Quantity)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Failed to add to cart"})
	}
	
	return c.JSON(200, gin.H{"data": cartItem})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	itemID := c.Param("item_id")
	id, err := strconv.ParseUint(itemID, 10, 32)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Invalid item ID"})
	}
	
	// Validate ID
	if err := validator.ValidateCartItemID(uint(id)); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	var req models.UpdateCartItemRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateUpdateCartItemRequest(req.Quantity); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Update
	cartItem, err := h.cartService.UpdateCartItem(uint(id), req.Quantity)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Failed to update cart item"})
	}
	
	return c.JSON(200, gin.H{"data": cartItem})
}
```

## 4. ORDER HANDLER EXAMPLE

```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateCreateOrderRequest(
		req.ShippingAddress, req.PaymentMethod,
	); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Get user ID
	userID, _ := c.Get("user_id")
	
	// Normalize
	req.ShippingAddress = validator.TrimWhitespace(req.ShippingAddress)
	
	// Create order
	order, err := h.orderService.Create(userID.(uint), &req)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Failed to create order"})
	}
	
	return c.JSON(201, gin.H{"data": order})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.ParseUint(orderID, 10, 32)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Invalid order ID"})
	}
	
	// Validate ID
	if err := validator.ValidateOrderID(uint(id)); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	var req models.UpdateOrderStatusRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateUpdateOrderStatusRequest(req.Status); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Get current order to check status transition
	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		return c.JSON(404, gin.H{"error": "Order not found"})
	}
	
	// Validate status transition
	if err := validator.ValidateOrderStatusTransition(order.Status, req.Status); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Update status
	updatedOrder, err := h.orderService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to update order status"})
	}
	
	return c.JSON(200, gin.H{"data": updatedOrder})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.ParseUint(orderID, 10, 32)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Invalid order ID"})
	}
	
	// Validate ID
	if err := validator.ValidateOrderID(uint(id)); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Get order
	order, err := h.orderService.GetByID(uint(id))
	if err != nil {
		return c.JSON(404, gin.H{"error": "Order not found"})
	}
	
	// Check if order can be cancelled
	if err := validator.ValidateCancelOrderRequest(order.Status); err != nil {
		// Only pending and processing orders can be cancelled
		if order.Status != "pending" && order.Status != "processing" {
			return c.JSON(400, gin.H{"error": "Cannot cancel order with status: " + order.Status})
		}
	}
	
	// Cancel order
	updatedOrder, err := h.orderService.Cancel(uint(id))
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to cancel order"})
	}
	
	return c.JSON(200, gin.H{"data": updatedOrder})
}
```

## 5. POST HANDLER EXAMPLE

```go
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateCreatePostRequest(req.Title, req.Content); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Get user ID
	userID, _ := c.Get("user_id")
	
	// Normalize
	req.Title = validator.TrimWhitespace(req.Title)
	req.Content = validator.TrimWhitespace(req.Content)
	
	// Create post
	post, err := h.postService.Create(userID.(uint), &req)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to create post"})
	}
	
	return c.JSON(201, gin.H{"data": post})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	id, err := strconv.ParseUint(postID, 10, 32)
	if err != nil {
		return c.JSON(400, gin.H{"error": "Invalid post ID"})
	}
	
	// Validate ID
	if err := validator.ValidatePostID(uint(id)); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	var req models.UpdatePostRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		return c.JSON(400, gin.H{"error": "Invalid request format"})
	}
	
	// Validate
	if err := validator.ValidateUpdatePostRequest(req.Title, req.Content, req.IsActive); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	
	// Update
	post, err := h.postService.Update(uint(id), &req)
	if err != nil {
		return c.JSON(500, gin.H{"error": "Failed to update post"})
	}
	
	return c.JSON(200, gin.H{"data": post})
}
```

## 6. ERROR RESPONSE HELPER

Create a common error response helper:

```go
// internal/api/handlers/error_helper.go
package handlers

import (
	"gin-ecommerce-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

func HandleValidationError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
			"status": "validation_error",
		})
		return true
	}
	return false
}

func HandleValidationErrors(c *gin.Context, errors validator.ValidationErrors) bool {
	if errors.HasErrors() {
		c.JSON(400, gin.H{
			"error": "Validation failed",
			"status": "validation_error",
			"details": errors.All(),
		})
		return true
	}
	return false
}
```

## Key Best Practices

1. **Always validate after binding**: Don't trust incoming JSON
2. **Normalize data**: Convert to consistent format (lowercase emails, trimmed strings)
3. **Check status transitions**: For entities with state machines (orders)
4. **Validate relationships**: Check that referenced IDs exist
5. **Return clear errors**: Include field names and what's wrong
6. **Use constants**: Define allowed statuses/methods in one place
7. **Test edge cases**: Empty values, nulls, boundary values
8. **Document requirements**: Keep README updated with validation rules

*/
