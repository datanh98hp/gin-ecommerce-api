package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars to restore later
	originalEnvVars := map[string]string{
		"SERVER_PORT":     os.Getenv("SERVER_PORT"),
		"ENV":             os.Getenv("ENV"),
		"DB_HOST":         os.Getenv("DB_HOST"),
		"DB_PORT":         os.Getenv("DB_PORT"),
		"DB_USER":         os.Getenv("DB_USER"),
		"DB_PASSWORD":     os.Getenv("DB_PASSWORD"),
		"DB_NAME":         os.Getenv("DB_NAME"),
		"DB_SSLMODE":      os.Getenv("DB_SSLMODE"),
		"JWT_SECRET":      os.Getenv("JWT_SECRET"),
		"JWT_EXPIRE_TIME": os.Getenv("JWT_EXPIRE_TIME"),
	}

	// Cleanup function to restore env vars
	defer func() {
		for key, val := range originalEnvVars {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}
	}()

	t.Run("Load with default values", func(t *testing.T) {
		// Clear all env vars
		for key := range originalEnvVars {
			os.Unsetenv(key)
		}

		cfg := Load()

		if cfg.Server.Port != "8080" {
			t.Errorf("Expected default port '8080', got '%s'", cfg.Server.Port)
		}
		if cfg.Server.Env != "development" {
			t.Errorf("Expected default env 'development', got '%s'", cfg.Server.Env)
		}
		if cfg.Database.Host != "localhost" {
			t.Errorf("Expected default DB host 'localhost', got '%s'", cfg.Database.Host)
		}
		if cfg.JWT.ExpireTime != 24 {
			t.Errorf("Expected default JWT expire time 24, got %d", cfg.JWT.ExpireTime)
		}
	})

	t.Run("Load with custom values", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "3000")
		os.Setenv("ENV", "production")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_USER", "customuser")
		os.Setenv("DB_PASSWORD", "custompass")
		os.Setenv("DB_NAME", "customdb")
		os.Setenv("DB_SSLMODE", "require")
		os.Setenv("JWT_SECRET", "custom-secret")
		os.Setenv("JWT_EXPIRE_TIME", "48")

		cfg := Load()

		if cfg.Server.Port != "3000" {
			t.Errorf("Expected port '3000', got '%s'", cfg.Server.Port)
		}
		if cfg.Server.Env != "production" {
			t.Errorf("Expected env 'production', got '%s'", cfg.Server.Env)
		}
		if cfg.Database.Host != "db.example.com" {
			t.Errorf("Expected DB host 'db.example.com', got '%s'", cfg.Database.Host)
		}
		if cfg.Database.Port != "5433" {
			t.Errorf("Expected DB port '5433', got '%s'", cfg.Database.Port)
		}
		if cfg.Database.User != "customuser" {
			t.Errorf("Expected DB user 'customuser', got '%s'", cfg.Database.User)
		}
		if cfg.JWT.Secret != "custom-secret" {
			t.Errorf("Expected JWT secret 'custom-secret', got '%s'", cfg.JWT.Secret)
		}
		if cfg.JWT.ExpireTime != 48 {
			t.Errorf("Expected JWT expire time 48, got %d", cfg.JWT.ExpireTime)
		}
	})

	t.Run("Load with partial custom values", func(t *testing.T) {
		// Clear all env vars
		for key := range originalEnvVars {
			os.Unsetenv(key)
		}

		os.Setenv("SERVER_PORT", "9000")
		os.Setenv("DB_NAME", "myapp")

		cfg := Load()

		// Custom values
		if cfg.Server.Port != "9000" {
			t.Errorf("Expected port '9000', got '%s'", cfg.Server.Port)
		}
		if cfg.Database.DBName != "myapp" {
			t.Errorf("Expected DB name 'myapp', got '%s'", cfg.Database.DBName)
		}

		// Default values for unset vars
		if cfg.Server.Env != "development" {
			t.Errorf("Expected default env 'development', got '%s'", cfg.Server.Env)
		}
		if cfg.Database.Host != "localhost" {
			t.Errorf("Expected default DB host 'localhost', got '%s'", cfg.Database.Host)
		}
	})

	t.Run("Invalid JWT_EXPIRE_TIME falls back to default", func(t *testing.T) {
		os.Setenv("JWT_EXPIRE_TIME", "invalid")

		cfg := Load()

		if cfg.JWT.ExpireTime != 24 {
			t.Errorf("Expected default JWT expire time 24 for invalid input, got %d", cfg.JWT.ExpireTime)
		}
	})
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{
			name: "Development environment",
			env:  "development",
			want: true,
		},
		{
			name: "Production environment",
			env:  "production",
			want: false,
		},
		{
			name: "Staging environment",
			env:  "staging",
			want: false,
		},
		{
			name: "Empty environment",
			env:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}

			if got := cfg.IsDevelopment(); got != tt.want {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want bool
	}{
		{
			name: "Production environment",
			env:  "production",
			want: true,
		},
		{
			name: "Development environment",
			env:  "development",
			want: false,
		},
		{
			name: "Staging environment",
			env:  "staging",
			want: false,
		},
		{
			name: "Empty environment",
			env:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}

			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "Existing environment variable",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "Non-existing environment variable",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
		{
			name:         "Empty environment variable",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		want         int
	}{
		{
			name:         "Valid integer",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "42",
			want:         42,
		},
		{
			name:         "Invalid integer",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "invalid",
			want:         10,
		},
		{
			name:         "Empty value",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "",
			want:         10,
		},
		{
			name:         "Negative integer",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "-5",
			want:         -5,
		},
		{
			name:         "Zero",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "0",
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getEnvAsInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigStructure(t *testing.T) {
	cfg := Load()

	// Test that all nested structs are properly initialized
	if cfg.Server.Port == "" {
		t.Error("Server.Port should not be empty")
	}
	if cfg.Server.Env == "" {
		t.Error("Server.Env should not be empty")
	}
	if cfg.Database.Host == "" {
		t.Error("Database.Host should not be empty")
	}
	if cfg.Database.Port == "" {
		t.Error("Database.Port should not be empty")
	}
	if cfg.JWT.Secret == "" {
		t.Error("JWT.Secret should not be empty")
	}
	if cfg.JWT.ExpireTime == 0 {
		t.Error("JWT.ExpireTime should not be zero")
	}
}
