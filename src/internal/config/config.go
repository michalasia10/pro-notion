package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Database struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:""`
	Database string `env:"DATABASE" envDefault:"pro_notion"`
	Schema   string `env:"SCHEMA" envDefault:"public"`
}

func (d *Database) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.Schema,
	)
}

type Redis struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"6379"`
	Password string `env:"PASSWORD" envDefault:""`
}

func (r *Redis) URL() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type Notion struct {
	ClientID      string `env:"CLIENT_ID" envDefault:""`
	ClientSecret  string `env:"CLIENT_SECRET" envDefault:""`
	RedirectURL   string `env:"REDIRECT_URL" envDefault:"http://localhost:8080/api/v1/auth/notion/callback"`
	APIVersion    string `env:"API_VERSION" envDefault:"2022-06-28"`
	WebhookSecret string `env:"WEBHOOK_SECRET" envDefault:""`
}

type JWT struct {
	Secret string `env:"SECRET" envDefault:"your-secret-key"`
}

type EventBus struct {
	Transport            string `env:"TRANSPORT"`
	ConsumerGroup        string `env:"CONSUMER_GROUP" envDefault:"worker_group"`
	TimeoutIntervalMin   int    `env:"TIMEOUT_INTERNAL_MINUTE" envDefault:"5"`
	ClaimIntervalSeconds int    `env:"CLAIM_INTERAL_SECONDS" envDefault:"5"`
}

func (eb *EventBus) GetTimeoutInterval() time.Duration {
	return time.Duration(eb.TimeoutIntervalMin) * time.Minute
}

func (eb *EventBus) GetClaimInterval() time.Duration {
	return time.Duration(eb.ClaimIntervalSeconds) * time.Second
}

type Async struct {
	Concurrency int    `env:"CONCURRENCY" envDefault:"10"`
	Queues      IntMap `env:"QUEUES" envDefault:"{\"critical\":6,\"default\":3}"`
}

// Config holds all configuration for the application
type Config struct {
	Port int `env:"PORT" envDefault:"8080"`

	Database Database `envPrefix:"BLUEPRINT_DB_"`
	Redis    Redis    `envPrefix:"REDIS_"`
	Notion   Notion   `envPrefix:"NOTION_"`
	JWT      JWT      `envPrefix:"JWT_"`
	EventBus EventBus `envPrefix:"EVENT_BUS_"`
	Async    Async    `envPrefix:"ASYNC_"`
}

var cfg *Config

func Load() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// Validate required config
	if cfg.Notion.ClientID == "" || cfg.Notion.ClientSecret == "" {
		log.Println("Warning: Notion Client ID and Secret not configured. OAuth flow will not work.")
	}

	return cfg
}

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

func (c *Config) DatabaseURL() string {
	return c.Database.DatabaseURL()
}

func (c *Config) RedisURL() string {
	return c.Redis.URL()
}
