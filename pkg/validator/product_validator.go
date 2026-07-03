package validator

import "fmt"

// ProductValidator provides validation for product-related requests

// ValidateCreateProductRequest validates product creation request
func ValidateCreateProductRequest(name string, description string, price float64, stock int, category, typeField, imageURL string) error {
	// Validate name
	if len(name) == 0 {
		return fmt.Errorf("product name is required")
	}

	name = TrimWhitespace(name)
	if len(name) < 3 {
		return fmt.Errorf("product name must be at least 3 characters long")
	}

	if len(name) > 255 {
		return fmt.Errorf("product name must not exceed 255 characters")
	}

	// Validate description (optional)
	if len(description) > 0 {
		if len(description) < 10 {
			return fmt.Errorf("product description must be at least 10 characters long if provided")
		}
		if len(description) > 2000 {
			return fmt.Errorf("product description must not exceed 2000 characters")
		}
	}

	// Validate price
	if price <= 0 {
		return fmt.Errorf("product price must be greater than 0")
	}

	if price > 999999.99 {
		return fmt.Errorf("product price must not exceed 999999.99")
	}

	// Validate stock
	if stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}

	if stock > 1000000 {
		return fmt.Errorf("product stock must not exceed 1000000")
	}

	// Validate category (optional)
	if len(category) > 0 {
		if len(category) > 100 {
			return fmt.Errorf("category must not exceed 100 characters")
		}
	}

	// Validate type (optional)
	if len(typeField) > 0 {
		if len(typeField) > 100 {
			return fmt.Errorf("type must not exceed 100 characters")
		}
	}

	// Validate image URL (optional)
	if len(imageURL) > 0 {
		if !IsValidURL(imageURL) {
			return fmt.Errorf("invalid image URL format")
		}
	}

	return nil
}

// ValidateUpdateProductRequest validates product update request
func ValidateUpdateProductRequest(name *string, description *string, price *float64, stock *int, category *string, typeField *string, imageURL *string, isActive *bool) error {
	if name != nil {
		if len(*name) == 0 {
			return fmt.Errorf("product name cannot be empty")
		}
		trimmedName := TrimWhitespace(*name)
		if len(trimmedName) < 3 {
			return fmt.Errorf("product name must be at least 3 characters long")
		}
		if len(trimmedName) > 255 {
			return fmt.Errorf("product name must not exceed 255 characters")
		}
	}

	if description != nil {
		if len(*description) > 0 {
			if len(*description) < 10 {
				return fmt.Errorf("product description must be at least 10 characters long")
			}
			if len(*description) > 2000 {
				return fmt.Errorf("product description must not exceed 2000 characters")
			}
		}
	}

	if price != nil {
		if *price <= 0 {
			return fmt.Errorf("product price must be greater than 0")
		}
		if *price > 999999.99 {
			return fmt.Errorf("product price must not exceed 999999.99")
		}
	}

	if stock != nil {
		if *stock < 0 {
			return fmt.Errorf("product stock cannot be negative")
		}
		if *stock > 1000000 {
			return fmt.Errorf("product stock must not exceed 1000000")
		}
	}

	if category != nil {
		if len(*category) > 100 {
			return fmt.Errorf("category must not exceed 100 characters")
		}
	}

	if typeField != nil {
		if len(*typeField) > 100 {
			return fmt.Errorf("type must not exceed 100 characters")
		}
	}

	if imageURL != nil {
		if len(*imageURL) > 0 {
			if !IsValidURL(*imageURL) {
				return fmt.Errorf("invalid image URL format")
			}
		}
	}

	return nil
}

// ValidateProductID validates product ID
func ValidateProductID(productID uint) error {
	if productID == 0 {
		return fmt.Errorf("product ID is required and must be greater than 0")
	}
	return nil
}
