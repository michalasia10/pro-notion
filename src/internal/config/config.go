package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port int

	// Database configuration
	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
		Schema   string
	}

	// Redis configuration
	Redis struct {
		Host     string
		Port     string
		Password string
	}

	// Notion API configuration
	Notion struct {
		ClientID      string
		ClientSecret  string
		RedirectURL   string
		APIVersion    string
		WebhookSecret string
	}

	// JWT configuration
	JWT struct {
		Secret string
	}
	Async struct {
		Concurrency int
		Queues      map[string]int
	}
}

var cfg *Config

// Load initializes the configuration from environment variables
func Load() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{}

	// Server
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}
	cfg.Port = port

	// Database
	cfg.Database.Host = getEnv("BLUEPRINT_DB_HOST", "localhost")
	cfg.Database.Port = getEnv("BLUEPRINT_DB_PORT", "5432")
	cfg.Database.Username = getEnv("BLUEPRINT_DB_USERNAME", "postgres")
	cfg.Database.Password = getEnv("BLUEPRINT_DB_PASSWORD", "")
	cfg.Database.Database = getEnv("BLUEPRINT_DB_DATABASE", "pro_notion")
	cfg.Database.Schema = getEnv("BLUEPRINT_DB_SCHEMA", "public")

	// Redis
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")

	// Notion
	cfg.Notion.ClientID = getEnv("NOTION_CLIENT_ID", "")
	cfg.Notion.ClientSecret = getEnv("NOTION_CLIENT_SECRET", "")
	cfg.Notion.RedirectURL = getEnv("NOTION_REDIRECT_URL", "http://localhost:8080/api/v1/auth/notion/callback")
	cfg.Notion.APIVersion = getEnv("NOTION_API_VERSION", "2022-06-28")
	cfg.Notion.WebhookSecret = getEnv("NOTION_WEBHOOK_SECRET", "")

	// JWT
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")

	// Validate required config
	if cfg.Notion.ClientID == "" || cfg.Notion.ClientSecret == "" {
		log.Println("Warning: Notion Client ID and Secret not configured. OAuth flow will not work.")
	}

	// Async
	asyncConcurrency, err := strconv.Atoi(getEnv("ASYNC_CONCURRENCY", "10"))
	if err != nil {
		log.Fatalf("Invalid ASYNC_CONCURRENCY value: %v", err)
	}
	cfg.Async.Concurrency = asyncConcurrency

	asyncQueues := getEnv("ASYNC_QUEUES", "{\"critical\": 6, \"default\": 3}")
	cfg.Async.Queues = make(map[string]int)
	err = json.Unmarshal([]byte(asyncQueues), &cfg.Async.Queues)
	if err != nil {
		log.Fatalf("Invalid ASYNC_QUEUES value: %v", err)
	}

	return cfg
}

// Get returns the loaded configuration
func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

// SetForTests allows setting config directly for testing purposes
func SetForTests(testCfg *Config) {
	cfg = testCfg
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// DatabaseURL returns formatted database connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Schema,
	)
}

// RedisURL returns formatted redis connection string
func (c *Config) RedisURL() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
