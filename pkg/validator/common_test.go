package validator

import "testing"

// Tests for common validators

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@example.co.uk", true},
		{"invalid@", false},
		{"@example.com", false},
		{"plainaddress", false},
		{"", false},
		{string(make([]byte, 255)), false}, // Too long
	}

	for _, tt := range tests {
		if IsValidEmail(tt.email) != tt.valid {
			t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, IsValidEmail(tt.email), tt.valid)
		}
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"valid_user", true},
		{"user-123", true},
		{"abc", true},
		{"ab", false},                     // Too short
		{"user name", false},              // Contains space
		{"", false},                       // Empty
		{string(make([]byte, 31)), false}, // Too long
	}

	for _, tt := range tests {
		if IsValidUsername(tt.username) != tt.valid {
			t.Errorf("IsValidUsername(%q) = %v, want %v", tt.username, IsValidUsername(tt.username), tt.valid)
		}
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"password123", true},
		{"short", false}, // Too short
		{"validpass", true},
		{"", false},
	}

	for _, tt := range tests {
		if IsValidPassword(tt.password) != tt.valid {
			t.Errorf("IsValidPassword(%q) = %v, want %v", tt.password, IsValidPassword(tt.password), tt.valid)
		}
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"John", true},
		{"Mary-Jane", true},
		{"O'Brien", true},
		{"J", false},       // Too short
		{"John123", false}, // Contains number
		{"", false},        // Empty
	}

	for _, tt := range tests {
		if IsValidName(tt.name) != tt.valid {
			t.Errorf("IsValidName(%q) = %v, want %v", tt.name, IsValidName(tt.name), tt.valid)
		}
	}
}

func TestIsValidPrice(t *testing.T) {
	tests := []struct {
		price float64
		valid bool
	}{
		{10.99, true},
		{0.01, true},
		{0, false},
		{-10.99, false},
	}

	for _, tt := range tests {
		if IsValidPrice(tt.price) != tt.valid {
			t.Errorf("IsValidPrice(%v) = %v, want %v", tt.price, IsValidPrice(tt.price), tt.valid)
		}
	}
}

func TestIsValidStock(t *testing.T) {
	tests := []struct {
		stock int
		valid bool
	}{
		{0, true},
		{100, true},
		{-1, false},
	}

	for _, tt := range tests {
		if IsValidStock(tt.stock) != tt.valid {
			t.Errorf("IsValidStock(%v) = %v, want %v", tt.stock, IsValidStock(tt.stock), tt.valid)
		}
	}
}

func TestIsValidQuantity(t *testing.T) {
	tests := []struct {
		quantity int
		valid    bool
	}{
		{1, true},
		{100, true},
		{0, false},
		{-1, false},
	}

	for _, tt := range tests {
		if IsValidQuantity(tt.quantity) != tt.valid {
			t.Errorf("IsValidQuantity(%v) = %v, want %v", tt.quantity, IsValidQuantity(tt.quantity), tt.valid)
		}
	}
}

func TestIsValidOrderStatus(t *testing.T) {
	tests := []struct {
		status string
		valid  bool
	}{
		{"pending", true},
		{"processing", true},
		{"shipped", true},
		{"delivered", true},
		{"cancelled", true},
		{"PENDING", true}, // Case insensitive
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		if IsValidOrderStatus(tt.status) != tt.valid {
			t.Errorf("IsValidOrderStatus(%q) = %v, want %v", tt.status, IsValidOrderStatus(tt.status), tt.valid)
		}
	}
}

func TestIsValidPaymentMethod(t *testing.T) {
	tests := []struct {
		method string
		valid  bool
	}{
		{"credit_card", true},
		{"debit_card", true},
		{"paypal", true},
		{"bank_transfer", true},
		{"cash", true},
		{"CREDIT_CARD", true}, // Case insensitive
		{"check", false},
		{"", false},
	}

	for _, tt := range tests {
		if IsValidPaymentMethod(tt.method) != tt.valid {
			t.Errorf("IsValidPaymentMethod(%q) = %v, want %v", tt.method, IsValidPaymentMethod(tt.method), tt.valid)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Test@Example.com", "test@example.com"},
		{"  user@example.com  ", "user@example.com"},
	}

	for _, tt := range tests {
		if result := NormalizeEmail(tt.input); result != tt.expected {
			t.Errorf("NormalizeEmail(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
