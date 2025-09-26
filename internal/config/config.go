// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		RESTPort        int           `yaml:"rest_port"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	MessageBroker struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		VHost    string `yaml:"vhost"`
	} `yaml:"message_broker"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	Auth struct {
		JWTSecret     string        `yaml:"jwt_secret"`
		JWTExpiry     time.Duration `yaml:"jwt_expiry"`
		RefreshSecret string        `yaml:"jwt_refresh_secret"`
		RefreshExpiry time.Duration `yaml:"jwt_refresh_expiry"`
	} `yaml:"auth"`

	OIDC struct {
		Enabled      bool     `yaml:"enabled"`
		ProviderURL  string   `yaml:"provider_url"`
		ClientID     string   `yaml:"client_id"`
		ClientSecret string   `yaml:"client_secret"`
		RedirectURL  string   `yaml:"redirect_url"`
		Scopes       []string `yaml:"scopes"`
	} `yaml:"oidc"`

	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`

	Metrics struct {
		Enabled bool   `yaml:"enabled"`
		Port    int    `yaml:"port"`
		Path    string `yaml:"path"`
	} `yaml:"metrics"`

	Migrations struct {
		Dir       string `yaml:"dir"`
		ConfigDir string `yaml:"config_dir"`
	} `yaml:"migrations"`

	NATSURL string `yaml:"nats_url"`

	// Notification configurations
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		From     string `yaml:"from"`
		TLS      bool   `yaml:"tls"`
	} `yaml:"smtp"`

	AT struct {
		APIKey   string `yaml:"api_key"`
		Username string `yaml:"username"`
		BaseURL  string `yaml:"base_url"`
	} `yaml:"at"`
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRabbitMQURL returns the RabbitMQ connection URL
func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.MessageBroker.User,
		c.MessageBroker.Password,
		c.MessageBroker.Host,
		c.MessageBroker.Port,
		c.MessageBroker.VHost,
	)
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Server: struct {
			RESTPort        int           `yaml:"rest_port"`
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
		}{
			RESTPort:        8080,
			ShutdownTimeout: 30 * time.Second,
		},

		Database: struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DBName   string `yaml:"dbname"`
			SSLMode  string `yaml:"sslmode"`
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "metabase"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "sil_backend_assessment_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},

		Redis: struct {
			Address  string `yaml:"address"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		}{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},

		Auth: struct {
			JWTSecret     string        `yaml:"jwt_secret"`
			JWTExpiry     time.Duration `yaml:"jwt_expiry"`
			RefreshSecret string        `yaml:"jwt_refresh_secret"`
			RefreshExpiry time.Duration `yaml:"jwt_refresh_expiry"`
		}{
			JWTSecret:     getEnv("JWT_SECRET", ""),
			JWTExpiry:     time.Duration(getEnvInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
			RefreshSecret: getEnv("REFRESH_SECRET", ""),
			RefreshExpiry: time.Duration(getEnvInt("REFRESH_EXPIRY_HOURS", 168)) * time.Hour,
		},

		OIDC: struct {
			Enabled      bool     `yaml:"enabled"`
			ProviderURL  string   `yaml:"provider_url"`
			ClientID     string   `yaml:"client_id"`
			ClientSecret string   `yaml:"client_secret"`
			RedirectURL  string   `yaml:"redirect_url"`
			Scopes       []string `yaml:"scopes"`
		}{
			Enabled:      getEnvBool("OIDC_ENABLED", false),
			ProviderURL:  getEnv("OIDC_PROVIDER_URL", "https://accounts.google.com"),
			ClientID:     getEnv("OIDC_CLIENT_ID", ""),
			ClientSecret: getEnv("OIDC_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/oidc/callback"),
			Scopes:       []string{"openid", "profile", "email"},
		},

		Logging: struct {
			Level string `yaml:"level"`
			File  string `yaml:"file"`
		}{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "audit.log"),
		},

		Metrics: struct {
			Enabled bool   `yaml:"enabled"`
			Port    int    `yaml:"port"`
			Path    string `yaml:"path"`
		}{
			Enabled: getEnvBool("METRICS_ENABLED", true),
			Port:    getEnvInt("METRICS_PORT", 8082),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},

		Migrations: struct {
			Dir       string `yaml:"dir"`
			ConfigDir string `yaml:"config_dir"`
		}{
			Dir:       getEnv("MIGRATION_DIR", "cmd/migrate/migrations"),
			ConfigDir: getEnv("CONFIG_DIR", "./config"),
		},

		NATSURL: getEnv("NATS_URL", "nats://localhost:4222"),
	}

	return config, nil
}

// Helper functions for environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
