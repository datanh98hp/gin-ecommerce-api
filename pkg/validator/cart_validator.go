package validator

import "fmt"

// CartValidator provides validation for cart-related requests

// ValidateAddToCartRequest validates add to cart request
func ValidateAddToCartRequest(productID uint, quantity int) error {
	if productID == 0 {
		return fmt.Errorf("product ID is required and must be greater than 0")
	}

	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if quantity > 1000 {
		return fmt.Errorf("quantity must not exceed 1000")
	}

	return nil
}

// ValidateUpdateCartItemRequest validates update cart item request
func ValidateUpdateCartItemRequest(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if quantity > 1000 {
		return fmt.Errorf("quantity must not exceed 1000")
	}

	return nil
}

// ValidateCartItemID validates cart item ID
func ValidateCartItemID(cartItemID uint) error {
	if cartItemID == 0 {
		return fmt.Errorf("cart item ID is required and must be greater than 0")
	}
	return nil
}

// ValidateCartID validates cart ID
func ValidateCartID(cartID uint) error {
	if cartID == 0 {
		return fmt.Errorf("cart ID is required and must be greater than 0")
	}
	return nil
}

// ValidateCartBulkOperation validates bulk cart operations
func ValidateCartBulkOperation(items []map[string]interface{}) error {
	if len(items) == 0 {
		return fmt.Errorf("cart items list cannot be empty")
	}

	if len(items) > 100 {
		return fmt.Errorf("cannot process more than 100 items at once")
	}

	totalQuantity := 0
	for i, item := range items {
		if productID, ok := item["product_id"]; !ok || productID == nil {
			return fmt.Errorf("product_id is required in item %d", i+1)
		}

		if quantity, ok := item["quantity"].(float64); ok {
			qty := int(quantity)
			if qty <= 0 {
				return fmt.Errorf("quantity must be greater than 0 in item %d", i+1)
			}
			if qty > 1000 {
				return fmt.Errorf("quantity must not exceed 1000 in item %d", i+1)
			}
			totalQuantity += qty
		} else {
			return fmt.Errorf("invalid quantity format in item %d", i+1)
		}
	}

	if totalQuantity > 5000 {
		return fmt.Errorf("total quantity across all items must not exceed 5000")
	}

	return nil
}
