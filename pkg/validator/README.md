# Validator Package Documentation

Complete validation package for all API endpoints in the gin-ecommerce-api project.

## Overview

The `pkg/validator` package provides comprehensive validation for all API request types. It includes:

- **Common validators**: Email, username, password, URL, phone number, and address validation
- **User validators**: Login, registration, password change, and profile update validation
- **Product validators**: Product creation, update, and product ID validation
- **Cart validators**: Cart operations, item quantity, and bulk operations
- **Order validators**: Order creation, status updates, and status transitions
- **Post validators**: Post creation, updates, and search
- **Pagination validators**: Page, sort order, date range validation
- **ID validators**: Generic ID validation

## File Structure

```
pkg/validator/
├── common.go                 # Core validation functions
├── user_validator.go         # User/auth validation
├── product_validator.go      # Product validation
├── cart_validator.go         # Cart validation
├── order_validator.go        # Order validation
├── post_validator.go         # Post validation
├── pagination_validator.go   # Pagination validation
├── id_validator.go          # ID validation
└── validator.go             # Main validator interface
```

## Usage Examples

### User Validation

```go
import "gin-ecommerce-api/pkg/validator"

// Validate login request
err := validator.ValidateLoginRequest(email, password)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate registration request
err := validator.ValidateRegisterRequest(email, username, password, firstName, lastName)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate password change
err := validator.ValidateChangePasswordRequest(oldPassword, newPassword, confirmPassword)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

### Product Validation

```go
// Validate product creation
err := validator.ValidateCreateProductRequest(name, description, price, stock, category, typeField, imageURL)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate product update
err := validator.ValidateUpdateProductRequest(&name, &desc, &price, &stock, &cat, &typ, &img, &active)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

### Cart Validation

```go
// Validate add to cart
err := validator.ValidateAddToCartRequest(productID, quantity)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate cart item update
err := validator.ValidateUpdateCartItemRequest(quantity)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

### Order Validation

```go
// Validate order creation
err := validator.ValidateCreateOrderRequest(shippingAddress, paymentMethod)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate order status update
err := validator.ValidateUpdateOrderStatusRequest(status)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate status transition
err := validator.ValidateOrderStatusTransition(currentStatus, newStatus)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

### Post Validation

```go
// Validate post creation
err := validator.ValidateCreatePostRequest(title, content)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate post update
err := validator.ValidateUpdatePostRequest(&title, &content, &isActive)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

### Pagination Validation

```go
// Validate basic pagination params
err := validator.ValidatePaginationParams(page, pageSize)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Validate and normalize with defaults
allowedFields := map[string]bool{
    "created_at": true,
    "name": true,
    "price": true,
}
params, err := validator.ValidateAndNormalizePaginationParams(page, pageSize, sortBy, sortOrder, allowedFields)
if err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}

// Use calculated values
offset := validator.CalculateOffset(params.Page, params.PageSize)
limit := validator.GetLimit(params.PageSize)
```

## Validation Rules

### Email
- Required
- Must be valid email format
- Maximum 254 characters

### Username
- Required
- 3-30 characters
- Alphanumeric, underscore, hyphen only

### Password
- Required
- 6-128 characters

### Names (First/Last)
- Required
- 2-50 characters
- Letters, spaces, hyphens, apostrophes only

### Product Name
- Required
- 3-255 characters
- Trimmed whitespace

### Product Description
- Optional
- If provided: 10-2000 characters

### Product Price
- Required
- Must be > 0
- Maximum 999999.99

### Product Stock
- Required
- Must be >= 0
- Maximum 1,000,000

### Quantity (Cart/Order)
- Required
- Must be > 0
- Maximum 1000 per item
- Total cart quantity: maximum 5000

### Shipping Address
- Required
- 5-255 characters
- Must be valid address format

### Payment Methods
- credit_card
- debit_card
- paypal
- bank_transfer
- cash

### Order Status
- pending
- processing
- shipped
- delivered
- cancelled

### Post Title
- Required
- 5-255 characters

### Post Content
- Required
- 20-10000 characters

### Pagination
- Page: 1-10000 (default: 1)
- Page Size: 1-100 (default: 10)
- Sort Order: asc or desc (default: desc)

## Common Helper Functions

### Text Normalization
```go
// Normalize email to lowercase
normalizedEmail := validator.NormalizeEmail(email)

// Normalize username to lowercase
normalizedUsername := validator.NormalizeUsername(username)

// Trim whitespace
trimmed := validator.TrimWhitespace(text)
```

### Format Validation
```go
// Check email format
isValid := validator.IsValidEmail(email)

// Check username format
isValid := validator.IsValidUsername(username)

// Check password strength
isValid := validator.IsValidPassword(password)

// Check URL format
isValid := validator.IsValidURL(url)

// Check phone number
isValid := validator.IsValidPhoneNumber(phone)

// Check address
isValid := validator.IsValidAddress(address)

// Check order status
isValid := validator.IsValidOrderStatus(status)

// Check payment method
isValid := validator.IsValidPaymentMethod(method)

// Check user role
isValid := validator.IsValidUserRole(role)
```

## Integration with Handlers

### Example Handler Integration

```go
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        return c.JSON(400, gin.H{"error": err.Error()})
    }
    
    // Validate request
    if err := validator.ValidateLoginRequest(req.Email, req.Password); err != nil {
        return c.JSON(400, gin.H{"error": err.Error()})
    }
    
    // Normalize for consistent processing
    req.Email = validator.NormalizeEmail(req.Email)
    
    // Continue with business logic
    user, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        return c.JSON(401, gin.H{"error": "Invalid credentials"})
    }
    
    return c.JSON(200, gin.H{"data": user})
}
```

## Error Handling Pattern

All validators return an `error` interface with descriptive messages:

```go
if err := validator.ValidateCreateProductRequest(name, desc, price, stock, cat, typ, img); err != nil {
    // err.Error() returns a descriptive message
    return c.JSON(400, gin.H{
        "error": err.Error(),
        "status": "validation_failed",
    })
}
```

## Best Practices

1. **Validate Early**: Call validators as soon as you receive data
2. **Normalize Data**: Use normalization functions before storing
3. **Use Type-Safe**: Leverage Go's type system with validators
4. **Custom Errors**: Wrap validation errors with context when needed
5. **Test Thoroughly**: Include validation in your test coverage

## Adding Custom Validators

To add custom validators:

1. Create a new file in `pkg/validator/`
2. Define validation functions following the pattern
3. Return descriptive error messages
4. Add tests for your validators

Example:
```go
// custom_validator.go
func ValidateCustomField(field string) error {
    if len(field) == 0 {
        return fmt.Errorf("custom field is required")
    }
    // Add more validation logic
    return nil
}
```

## Testing

All validators can be unit tested:

```go
func TestValidateLoginRequest(t *testing.T) {
    tests := []struct {
        email    string
        password string
        hasError bool
    }{
        {"test@example.com", "password123", false},
        {"invalid-email", "password123", true},
        {"test@example.com", "short", true},
    }
    
    for _, tt := range tests {
        err := validator.ValidateLoginRequest(tt.email, tt.password)
        if (err != nil) != tt.hasError {
            t.Errorf("got %v, want error %v", err, tt.hasError)
        }
    }
}
```
