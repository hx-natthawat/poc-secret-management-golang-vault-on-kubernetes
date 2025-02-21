package config

import "os"

// Config holds application configuration
type Config struct {
	VaultAddr     string
	VaultRoleName string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		VaultAddr:     getEnvOrDefault("VAULT_ADDR", "http://vault.vault:8200"),
		VaultRoleName: getEnvOrDefault("VAULT_ROLE_NAME", "app-role"),
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
