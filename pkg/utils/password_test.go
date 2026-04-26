package utils

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Short password",
			password: "pass",
			wantErr:  false,
		},
		{
			name:     "Long password within limit",
			password: strings.Repeat("a", 70),
			wantErr:  false,
		},
		{
			name:     "Special characters",
			password: "P@ssw0rd!#$%",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plaintext password")
				}
				// Verify it's a valid bcrypt hash
				if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") && !strings.HasPrefix(hash, "$2y$") {
					t.Error("HashPassword() did not return valid bcrypt hash")
				}
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "testpassword"
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to hash password: %v, %v", err1, err2)
	}

	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("HashPassword() should produce different hashes for the same password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		want           bool
	}{
		{
			name:           "Correct password",
			hashedPassword: hash,
			password:       password,
			want:           true,
		},
		{
			name:           "Incorrect password",
			hashedPassword: hash,
			password:       "wrongpassword",
			want:           false,
		},
		{
			name:           "Empty password",
			hashedPassword: hash,
			password:       "",
			want:           false,
		},
		{
			name:           "Case sensitive",
			hashedPassword: hash,
			password:       "TESTPASSWORD123",
			want:           false,
		},
		{
			name:           "Extra characters",
			hashedPassword: hash,
			password:       password + "extra",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.hashedPassword, tt.password)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	tests := []struct {
		name           string
		hashedPassword string
		password       string
		want           bool
	}{
		{
			name:           "Invalid hash format",
			hashedPassword: "not-a-valid-hash",
			password:       "password",
			want:           false,
		},
		{
			name:           "Empty hash",
			hashedPassword: "",
			password:       "password",
			want:           false,
		},
		{
			name:           "Plaintext as hash",
			hashedPassword: "password",
			password:       "password",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.hashedPassword, tt.password)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHashCost(t *testing.T) {
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Extract the cost from the hash
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("Failed to extract cost: %v", err)
	}

	// Verify it uses the default cost
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected cost %d, got %d", bcrypt.DefaultCost, cost)
	}
}

func TestPasswordRoundTrip(t *testing.T) {
	testPasswords := []string{
		"simple",
		"Complex123!@#",
		"with spaces in it",
		"unicode-密码-🔐",
		strings.Repeat("long", 15), // 60 bytes, within bcrypt's 72-byte limit
	}

	for _, password := range testPasswords {
		t.Run(password, func(t *testing.T) {
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() failed: %v", err)
			}

			if !CheckPassword(hash, password) {
				t.Error("CheckPassword() failed for correct password")
			}

			// Verify wrong passwords don't match
			if CheckPassword(hash, password+"wrong") {
				t.Error("CheckPassword() succeeded for incorrect password")
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkpassword"
	hash, _ := HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(hash, password)
	}
}
