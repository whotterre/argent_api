package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnvVarsWithoutDotEnvFile(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"PORT", "DATABASE_URL", "GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_REDIRECT_URL", "JWT_SECRET", "PAYSTACK_SECRET"}
	
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}
	
	// Clean up after test
	defer func() {
		for key, val := range originalEnv {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
	os.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/callback")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("PAYSTACK_SECRET", "test_paystack_secret")

	// Load config (no .env file should exist in test environment)
	cfg, err := LoadConfig()

	// Check that no error is returned even without .env file
	if err != nil {
		t.Fatalf("LoadConfig() returned an error when loading from env vars: %v", err)
	}

	// Verify that the config values are loaded from environment variables
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("Expected DatabaseURL to be 'postgres://user:pass@localhost:5432/testdb', got '%s'", cfg.DatabaseURL)
	}

	if cfg.GoogleClientID != "test_client_id" {
		t.Errorf("Expected GoogleClientID to be 'test_client_id', got '%s'", cfg.GoogleClientID)
	}

	if cfg.GoogleClientSecret != "test_client_secret" {
		t.Errorf("Expected GoogleClientSecret to be 'test_client_secret', got '%s'", cfg.GoogleClientSecret)
	}

	if cfg.GoogleRedirectURL != "http://localhost:8080/callback" {
		t.Errorf("Expected GoogleRedirectURL to be 'http://localhost:8080/callback', got '%s'", cfg.GoogleRedirectURL)
	}

	if cfg.JWTSecret != "test_jwt_secret" {
		t.Errorf("Expected JWTSecret to be 'test_jwt_secret', got '%s'", cfg.JWTSecret)
	}

	if cfg.PaystackSecret != "test_paystack_secret" {
		t.Errorf("Expected PaystackSecret to be 'test_paystack_secret', got '%s'", cfg.PaystackSecret)
	}

	t.Log("✓ Config successfully loaded from environment variables without .env file")
}

