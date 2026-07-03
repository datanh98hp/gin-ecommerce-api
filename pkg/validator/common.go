package validator

import (
	"regexp"
	"strings"
)

// Common validation rules and utilities

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	const emailPattern = `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

// IsValidUsername validates username format (alphanumeric, underscore, hyphen, 3-30 chars)
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return pattern.MatchString(username)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	if len(password) > 128 {
		return false
	}
	return true
}

// IsValidName validates name format (2-50 chars, letters and spaces)
func IsValidName(name string) bool {
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	pattern := regexp.MustCompile(`^[a-zA-Z\s-']+$`)
	return pattern.MatchString(name)
}

// IsValidPrice validates price (must be positive)
func IsValidPrice(price float64) bool {
	return price > 0
}

// IsValidStock validates stock quantity (must be non-negative)
func IsValidStock(stock int) bool {
	return stock >= 0
}

// IsValidQuantity validates quantity (must be positive)
func IsValidQuantity(quantity int) bool {
	return quantity > 0
}

// IsValidURL validates URL format
func IsValidURL(url string) bool {
	if len(url) == 0 {
		return true // Optional field
	}
	if len(url) > 2048 {
		return false
	}
	pattern := regexp.MustCompile(`^https?://[^\s]+$`)
	return pattern.MatchString(url)
}

// IsValidPhoneNumber validates phone number (basic format)
func IsValidPhoneNumber(phone string) bool {
	if len(phone) == 0 {
		return true // Optional
	}
	pattern := regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)
	if !pattern.MatchString(phone) {
		return false
	}
	// Remove non-digit characters and check length
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	return len(digits) >= 10 && len(digits) <= 15
}

// IsValidAddress validates address format
func IsValidAddress(address string) bool {
	if len(address) < 5 || len(address) > 255 {
		return false
	}
	return true
}

// IsValidOrderStatus validates order status
func IsValidOrderStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}
	return validStatuses[strings.ToLower(status)]
}

// IsValidPaymentMethod validates payment method
func IsValidPaymentMethod(method string) bool {
	validMethods := map[string]bool{
		"credit_card":   true,
		"debit_card":    true,
		"paypal":        true,
		"bank_transfer": true,
		"cash":          true,
	}
	return validMethods[strings.ToLower(method)]
}

// IsValidUserRole validates user role
func IsValidUserRole(role string) bool {
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
	}
	return validRoles[strings.ToLower(role)]
}

// TrimWhitespace trims leading/trailing whitespace
func TrimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeEmail normalizes email to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeUsername normalizes username to lowercase
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
