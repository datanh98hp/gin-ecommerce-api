package validator

import "fmt"

// OrderValidator provides validation for order-related requests

// ValidateCreateOrderRequest validates order creation request
func ValidateCreateOrderRequest(shippingAddress, paymentMethod string) error {
	// Validate shipping address
	if len(shippingAddress) == 0 {
		return fmt.Errorf("shipping address is required")
	}

	shippingAddress = TrimWhitespace(shippingAddress)
	if !IsValidAddress(shippingAddress) {
		return fmt.Errorf("shipping address must be 5-255 characters long")
	}

	// Validate payment method
	if len(paymentMethod) == 0 {
		return fmt.Errorf("payment method is required")
	}

	if !IsValidPaymentMethod(paymentMethod) {
		return fmt.Errorf("invalid payment method: must be one of 'credit_card', 'debit_card', 'paypal', 'bank_transfer', 'cash'")
	}

	return nil
}

// ValidateUpdateOrderStatusRequest validates order status update request
func ValidateUpdateOrderStatusRequest(status string) error {
	if len(status) == 0 {
		return fmt.Errorf("status is required")
	}

	if !IsValidOrderStatus(status) {
		return fmt.Errorf("invalid status: must be one of 'pending', 'processing', 'shipped', 'delivered', 'cancelled'")
	}

	return nil
}

// ValidateOrderID validates order ID
func ValidateOrderID(orderID uint) error {
	if orderID == 0 {
		return fmt.Errorf("order ID is required and must be greater than 0")
	}
	return nil
}

// ValidateCancelOrderRequest validates cancel order request
func ValidateCancelOrderRequest(status string) error {
	if len(status) == 0 {
		status = "cancelled"
	}

	if status != "cancelled" {
		return fmt.Errorf("order cancellation must set status to 'cancelled'")
	}

	return nil
}

// ValidateOrderTotalAmount validates order total amount
func ValidateOrderTotalAmount(totalAmount float64) error {
	if totalAmount <= 0 {
		return fmt.Errorf("order total amount must be greater than 0")
	}

	if totalAmount > 999999.99 {
		return fmt.Errorf("order total amount must not exceed 999999.99")
	}

	return nil
}

// ValidateOrderStatusTransition validates if status transition is allowed
func ValidateOrderStatusTransition(currentStatus, newStatus string) error {
	// Define valid transitions
	validTransitions := map[string]map[string]bool{
		"pending": {
			"processing": true,
			"cancelled":  true,
		},
		"processing": {
			"shipped":   true,
			"cancelled": true,
		},
		"shipped": {
			"delivered": true,
		},
		"delivered": {},
		"cancelled": {},
	}

	transitions, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	if !transitions[newStatus] {
		return fmt.Errorf("cannot transition from '%s' to '%s'", currentStatus, newStatus)
	}

	return nil
}

// ValidateOrderItem validates order item data
func ValidateOrderItem(productID uint, quantity int, price float64) error {
	if productID == 0 {
		return fmt.Errorf("product ID is required and must be greater than 0")
	}

	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if quantity > 1000 {
		return fmt.Errorf("quantity must not exceed 1000")
	}

	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if price > 999999.99 {
		return fmt.Errorf("price must not exceed 999999.99")
	}

	return nil
}
