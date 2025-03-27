package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config structure to hold DB settings
type Config struct {
	DBDriver           string `yaml:"db_driver"`
	DBSource           string `yaml:"db_source"`
	RedisAddr          string `yaml:"db_redis_addr"`
	DaisiApiUrl        string `yaml:"daisi_api_url"`
	DaisiApiSenderName string `yaml:"daisi_api_sender_name"`
	DaisiApiToken      string `yaml:"daisi_api_token"`
}

// CONFIG is the global configuration instance
var CONFIG *Config

// LoadConfig reads YAML based on the GO_ENV variable
func LoadConfig() Config {
	env := os.Getenv("GO_ENV") // Get current environment (development, staging, production)
	if env == "" {
		env = "development" // Default to development
	}

	// Read YAML file
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse YAML
	var fullConfig map[string]Config
	err = yaml.Unmarshal(yamlFile, &fullConfig)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	// Get config for the current environment
	config, exists := fullConfig[env]
	if !exists {
		log.Fatalf("No configuration found for environment: %s", env)
	}

	// ✅ Assign to global CONFIG pointer
	CONFIG = &config

	log.Println("🚀 Application running with loaded configuration for:", env)
	return config
}
