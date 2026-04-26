package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		email       string
		role        string
		secret      string
		expireHours int
		wantErr     bool
	}{
		{
			name:        "Valid token generation",
			userID:      1,
			email:       "test@example.com",
			role:        "user",
			secret:      "test-secret",
			expireHours: 24,
			wantErr:     false,
		},
		{
			name:        "Admin token generation",
			userID:      2,
			email:       "admin@example.com",
			role:        "admin",
			secret:      "test-secret",
			expireHours: 24,
			wantErr:     false,
		},
		{
			name:        "Short expiration",
			userID:      3,
			email:       "user@example.com",
			role:        "user",
			secret:      "test-secret",
			expireHours: 1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.role, tt.secret, tt.expireHours)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("GenerateToken() returned empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	userID := uint(1)
	email := "test@example.com"
	role := "user"

	// Generate a valid token
	validToken, err := GenerateToken(userID, email, role, secret, 24)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secret    string
		wantErr   bool
		wantID    uint
		wantEmail string
		wantRole  string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			secret:    secret,
			wantErr:   false,
			wantID:    userID,
			wantEmail: email,
			wantRole:  role,
		},
		{
			name:    "Invalid secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "invalid.token.here",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims.UserID != tt.wantID {
					t.Errorf("ValidateToken() userID = %v, want %v", claims.UserID, tt.wantID)
				}
				if claims.Email != tt.wantEmail {
					t.Errorf("ValidateToken() email = %v, want %v", claims.Email, tt.wantEmail)
				}
				if claims.Role != tt.wantRole {
					t.Errorf("ValidateToken() role = %v, want %v", claims.Role, tt.wantRole)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test-secret"
	userID := uint(1)
	email := "test@example.com"
	role := "user"

	// Generate a token that expires in a very short time
	token, err := GenerateToken(userID, email, role, secret, -1) // Expired token
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// Try to validate the expired token
	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Error("ValidateToken() should fail for expired token")
	}
}

func TestTokenClaims(t *testing.T) {
	secret := "test-secret"
	userID := uint(123)
	email := "testuser@example.com"
	role := "admin"
	expireHours := 48

	token, err := GenerateToken(userID, email, role, secret, expireHours)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Verify all claims
	if claims.UserID != userID {
		t.Errorf("Expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}

	// Check expiration time is approximately correct (within 1 minute tolerance)
	expectedExpiry := time.Now().Add(time.Duration(expireHours) * time.Hour)
	actualExpiry := claims.ExpiresAt.Time
	diff := actualExpiry.Sub(expectedExpiry)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("Token expiration time is off by %v", diff)
	}
}

func TestInvalidSigningMethod(t *testing.T) {
	// Create a token with a different signing method
	claims := &Claims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Use RSA instead of HMAC (our code expects HMAC)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // Using HS512 instead of HS256
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// This should still work as it's still HMAC-based
	_, err = ValidateToken(tokenString, "test-secret")
	if err != nil {
		t.Errorf("ValidateToken() should accept HS512: %v", err)
	}
}
