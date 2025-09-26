package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("Load valid config file", func(t *testing.T) {
		// Create a temporary config file
		configContent := `
server:
  rest_port: 8080
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  user: testuser
  password: testpass
  dbname: testdb
  sslmode: disable

redis:
  address: localhost:6379
  password: redispass
  db: 0

auth:
  jwt_secret: test-jwt-secret
  jwt_expiry: 24h
  jwt_refresh_secret: test-refresh-secret
  jwt_refresh_expiry: 168h

oidc:
  enabled: true
  provider_url: https://accounts.google.com
  client_id: test-client-id
  client_secret: test-client-secret
  redirect_url: http://localhost:8080/callback
  scopes: ["openid", "profile", "email"]

logging:
  level: debug
  file: test.log

metrics:
  enabled: true
  port: 8082
  path: /metrics

migrations:
  dir: migrations
  config_dir: config

nats_url: nats://localhost:4222

smtp:
  host: smtp.gmail.com
  port: 587
  username: test@example.com
  password: testpass
  from: noreply@example.com
  tls: true

at:
  api_key: test-api-key
  username: test-username
  base_url: https://api.africastalking.com
`

		// Create temporary file
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Load config
		cfg, err := Load(configFile)
		require.NoError(t, err)
		assert.NotNil(t, cfg)

		// Verify server config
		assert.Equal(t, 8080, cfg.Server.RESTPort)
		assert.Equal(t, 30*time.Second, cfg.Server.ShutdownTimeout)

		// Verify database config
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "testuser", cfg.Database.User)
		assert.Equal(t, "testpass", cfg.Database.Password)
		assert.Equal(t, "testdb", cfg.Database.DBName)
		assert.Equal(t, "disable", cfg.Database.SSLMode)

		// Verify Redis config
		assert.Equal(t, "localhost:6379", cfg.Redis.Address)
		assert.Equal(t, "redispass", cfg.Redis.Password)
		assert.Equal(t, 0, cfg.Redis.DB)

		// Verify auth config
		assert.Equal(t, "test-jwt-secret", cfg.Auth.JWTSecret)
		assert.Equal(t, 24*time.Hour, cfg.Auth.JWTExpiry)
		assert.Equal(t, "test-refresh-secret", cfg.Auth.RefreshSecret)
		assert.Equal(t, 168*time.Hour, cfg.Auth.RefreshExpiry)

		// Verify OIDC config
		assert.True(t, cfg.OIDC.Enabled)
		assert.Equal(t, "https://accounts.google.com", cfg.OIDC.ProviderURL)
		assert.Equal(t, "test-client-id", cfg.OIDC.ClientID)
		assert.Equal(t, "test-client-secret", cfg.OIDC.ClientSecret)
		assert.Equal(t, "http://localhost:8080/callback", cfg.OIDC.RedirectURL)
		assert.Equal(t, []string{"openid", "profile", "email"}, cfg.OIDC.Scopes)

		// Verify logging config
		assert.Equal(t, "debug", cfg.Logging.Level)
		assert.Equal(t, "test.log", cfg.Logging.File)

		// Verify metrics config
		assert.True(t, cfg.Metrics.Enabled)
		assert.Equal(t, 8082, cfg.Metrics.Port)
		assert.Equal(t, "/metrics", cfg.Metrics.Path)

		// Verify migrations config
		assert.Equal(t, "migrations", cfg.Migrations.Dir)
		assert.Equal(t, "config", cfg.Migrations.ConfigDir)

		// Verify NATS URL
		assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)

		// Verify SMTP config
		assert.Equal(t, "smtp.gmail.com", cfg.SMTP.Host)
		assert.Equal(t, 587, cfg.SMTP.Port)
		assert.Equal(t, "test@example.com", cfg.SMTP.Username)
		assert.Equal(t, "testpass", cfg.SMTP.Password)
		assert.Equal(t, "noreply@example.com", cfg.SMTP.From)
		assert.True(t, cfg.SMTP.TLS)

		// Verify AT config
		assert.Equal(t, "test-api-key", cfg.AT.APIKey)
		assert.Equal(t, "test-username", cfg.AT.Username)
		assert.Equal(t, "https://api.africastalking.com", cfg.AT.BaseURL)
	})

	t.Run("Load non-existent config file", func(t *testing.T) {
		cfg, err := Load("non-existent-file.yaml")

		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to read config file")
	})

	t.Run("Load invalid YAML config file", func(t *testing.T) {
		// Create a temporary config file with invalid YAML
		configContent := `
server:
  rest_port: 8080
  shutdown_timeout: 30s
invalid: yaml: content: [unclosed
`

		// Create temporary file
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "invalid-config.yaml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Load config
		cfg, err := Load(configFile)

		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse config file")
	})
}

func TestConfig_GetDSN(t *testing.T) {
	cfg := &Config{
		Database: struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DBName   string `yaml:"dbname"`
			SSLMode  string `yaml:"sslmode"`
		}{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	actualDSN := cfg.GetDSN()

	assert.Equal(t, expectedDSN, actualDSN)
}

func TestConfig_GetRabbitMQURL(t *testing.T) {
	cfg := &Config{
		MessageBroker: struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			VHost    string `yaml:"vhost"`
		}{
			Host:     "localhost",
			Port:     5672,
			User:     "guest",
			Password: "guest",
			VHost:    "/",
		},
	}

	expectedURL := "amqp://guest:guest@localhost:5672//"
	actualURL := cfg.GetRabbitMQURL()

	assert.Equal(t, expectedURL, actualURL)
}

func TestLoadFromEnv(t *testing.T) {
	t.Run("Load config with default values", func(t *testing.T) {
		// Clear environment variables
		envVars := []string{
			"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
			"REDIS_ADDRESS", "REDIS_PASSWORD", "REDIS_DB",
			"JWT_SECRET", "JWT_EXPIRY_HOURS", "REFRESH_SECRET", "REFRESH_EXPIRY_HOURS",
			"OIDC_ENABLED", "OIDC_PROVIDER_URL", "OIDC_CLIENT_ID", "OIDC_CLIENT_SECRET", "OIDC_REDIRECT_URL",
			"LOG_LEVEL", "LOG_FILE",
			"METRICS_ENABLED", "METRICS_PORT", "METRICS_PATH",
			"MIGRATION_DIR", "CONFIG_DIR",
			"NATS_URL",
		}

		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		cfg, err := LoadFromEnv()
		require.NoError(t, err)
		assert.NotNil(t, cfg)

		// Verify default values
		assert.Equal(t, 8080, cfg.Server.RESTPort)
		assert.Equal(t, 30*time.Second, cfg.Server.ShutdownTimeout)
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "metabase", cfg.Database.User)
		assert.Equal(t, "ffdd1f0a568f407daaa2e176b5fd5481", cfg.Database.Password)
		assert.Equal(t, "sil_backend_assessment_db", cfg.Database.DBName)
		assert.Equal(t, "disable", cfg.Database.SSLMode)
		assert.Equal(t, "localhost:6379", cfg.Redis.Address)
		assert.Equal(t, "ffdd1f0a568f407daaa2e176b5fd5481", cfg.Redis.Password)
		assert.Equal(t, 0, cfg.Redis.DB)
		assert.Equal(t, "your-jwt-secret", cfg.Auth.JWTSecret)
		assert.Equal(t, 24*time.Hour, cfg.Auth.JWTExpiry)
		assert.Equal(t, "your-refresh-secret", cfg.Auth.RefreshSecret)
		assert.Equal(t, 168*time.Hour, cfg.Auth.RefreshExpiry)
		assert.False(t, cfg.OIDC.Enabled)
		assert.Equal(t, "https://accounts.google.com", cfg.OIDC.ProviderURL)
		assert.Equal(t, "", cfg.OIDC.ClientID)
		assert.Equal(t, "", cfg.OIDC.ClientSecret)
		assert.Equal(t, "http://localhost:8080/auth/oidc/callback", cfg.OIDC.RedirectURL)
		assert.Equal(t, []string{"openid", "profile", "email"}, cfg.OIDC.Scopes)
		assert.Equal(t, "info", cfg.Logging.Level)
		assert.Equal(t, "audit.log", cfg.Logging.File)
		assert.True(t, cfg.Metrics.Enabled)
		assert.Equal(t, 8082, cfg.Metrics.Port)
		assert.Equal(t, "/metrics", cfg.Metrics.Path)
		assert.Equal(t, "cmd/migrate/migrations", cfg.Migrations.Dir)
		assert.Equal(t, "./config", cfg.Migrations.ConfigDir)
		assert.Equal(t, "nats://localhost:4222", cfg.NATSURL)
	})

	t.Run("Load config with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("DB_HOST", "test-host")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_USER", "test-user")
		os.Setenv("DB_PASSWORD", "test-password")
		os.Setenv("DB_NAME", "test-db")
		os.Setenv("DB_SSL_MODE", "require")
		os.Setenv("REDIS_ADDRESS", "redis-host:6380")
		os.Setenv("REDIS_PASSWORD", "redis-password")
		os.Setenv("REDIS_DB", "1")
		os.Setenv("JWT_SECRET", "custom-jwt-secret")
		os.Setenv("JWT_EXPIRY_HOURS", "12")
		os.Setenv("REFRESH_SECRET", "custom-refresh-secret")
		os.Setenv("REFRESH_EXPIRY_HOURS", "72")
		os.Setenv("OIDC_ENABLED", "true")
		os.Setenv("OIDC_PROVIDER_URL", "https://auth.example.com")
		os.Setenv("OIDC_CLIENT_ID", "custom-client-id")
		os.Setenv("OIDC_CLIENT_SECRET", "custom-client-secret")
		os.Setenv("OIDC_REDIRECT_URL", "https://app.example.com/callback")
		os.Setenv("LOG_LEVEL", "error")
		os.Setenv("LOG_FILE", "error.log")
		os.Setenv("METRICS_ENABLED", "false")
		os.Setenv("METRICS_PORT", "9090")
		os.Setenv("METRICS_PATH", "/custom-metrics")
		os.Setenv("MIGRATION_DIR", "custom-migrations")
		os.Setenv("CONFIG_DIR", "custom-config")
		os.Setenv("NATS_URL", "nats://custom-host:4222")

		cfg, err := LoadFromEnv()
		require.NoError(t, err)
		assert.NotNil(t, cfg)

		// Verify environment variable values
		assert.Equal(t, "test-host", cfg.Database.Host)
		assert.Equal(t, 3306, cfg.Database.Port)
		assert.Equal(t, "test-user", cfg.Database.User)
		assert.Equal(t, "test-password", cfg.Database.Password)
		assert.Equal(t, "test-db", cfg.Database.DBName)
		assert.Equal(t, "require", cfg.Database.SSLMode)
		assert.Equal(t, "redis-host:6380", cfg.Redis.Address)
		assert.Equal(t, "redis-password", cfg.Redis.Password)
		assert.Equal(t, 1, cfg.Redis.DB)
		assert.Equal(t, "custom-jwt-secret", cfg.Auth.JWTSecret)
		assert.Equal(t, 12*time.Hour, cfg.Auth.JWTExpiry)
		assert.Equal(t, "custom-refresh-secret", cfg.Auth.RefreshSecret)
		assert.Equal(t, 72*time.Hour, cfg.Auth.RefreshExpiry)
		assert.True(t, cfg.OIDC.Enabled)
		assert.Equal(t, "https://auth.example.com", cfg.OIDC.ProviderURL)
		assert.Equal(t, "custom-client-id", cfg.OIDC.ClientID)
		assert.Equal(t, "custom-client-secret", cfg.OIDC.ClientSecret)
		assert.Equal(t, "https://app.example.com/callback", cfg.OIDC.RedirectURL)
		assert.Equal(t, "error", cfg.Logging.Level)
		assert.Equal(t, "error.log", cfg.Logging.File)
		assert.False(t, cfg.Metrics.Enabled)
		assert.Equal(t, 9090, cfg.Metrics.Port)
		assert.Equal(t, "/custom-metrics", cfg.Metrics.Path)
		assert.Equal(t, "custom-migrations", cfg.Migrations.Dir)
		assert.Equal(t, "custom-config", cfg.Migrations.ConfigDir)
		assert.Equal(t, "nats://custom-host:4222", cfg.NATSURL)

		// Clean up environment variables
		for _, envVar := range []string{
			"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
			"REDIS_ADDRESS", "REDIS_PASSWORD", "REDIS_DB",
			"JWT_SECRET", "JWT_EXPIRY_HOURS", "REFRESH_SECRET", "REFRESH_EXPIRY_HOURS",
			"OIDC_ENABLED", "OIDC_PROVIDER_URL", "OIDC_CLIENT_ID", "OIDC_CLIENT_SECRET", "OIDC_REDIRECT_URL",
			"LOG_LEVEL", "LOG_FILE",
			"METRICS_ENABLED", "METRICS_PORT", "METRICS_PATH",
			"MIGRATION_DIR", "CONFIG_DIR",
			"NATS_URL",
		} {
			os.Unsetenv(envVar)
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("Get existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default-value")
		assert.Equal(t, "test-value", result)
	})

	t.Run("Get non-existing environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTING_VAR")

		result := getEnv("NON_EXISTING_VAR", "default-value")
		assert.Equal(t, "default-value", result)
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Run("Get valid integer environment variable", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		result := getEnvInt("TEST_INT", 0)
		assert.Equal(t, 42, result)
	})

	t.Run("Get invalid integer environment variable", func(t *testing.T) {
		os.Setenv("TEST_INVALID_INT", "not-a-number")
		defer os.Unsetenv("TEST_INVALID_INT")

		result := getEnvInt("TEST_INVALID_INT", 10)
		assert.Equal(t, 10, result)
	})

	t.Run("Get non-existing integer environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTING_INT")

		result := getEnvInt("NON_EXISTING_INT", 5)
		assert.Equal(t, 5, result)
	})
}

func TestGetEnvBool(t *testing.T) {
	t.Run("Get valid boolean environment variable - true", func(t *testing.T) {
		os.Setenv("TEST_BOOL_TRUE", "true")
		defer os.Unsetenv("TEST_BOOL_TRUE")

		result := getEnvBool("TEST_BOOL_TRUE", false)
		assert.True(t, result)
	})

	t.Run("Get valid boolean environment variable - false", func(t *testing.T) {
		os.Setenv("TEST_BOOL_FALSE", "false")
		defer os.Unsetenv("TEST_BOOL_FALSE")

		result := getEnvBool("TEST_BOOL_FALSE", true)
		assert.False(t, result)
	})

	t.Run("Get valid boolean environment variable - 1", func(t *testing.T) {
		os.Setenv("TEST_BOOL_1", "1")
		defer os.Unsetenv("TEST_BOOL_1")

		result := getEnvBool("TEST_BOOL_1", false)
		assert.True(t, result)
	})

	t.Run("Get valid boolean environment variable - 0", func(t *testing.T) {
		os.Setenv("TEST_BOOL_0", "0")
		defer os.Unsetenv("TEST_BOOL_0")

		result := getEnvBool("TEST_BOOL_0", true)
		assert.False(t, result)
	})

	t.Run("Get invalid boolean environment variable", func(t *testing.T) {
		os.Setenv("TEST_INVALID_BOOL", "maybe")
		defer os.Unsetenv("TEST_INVALID_BOOL")

		result := getEnvBool("TEST_INVALID_BOOL", true)
		assert.True(t, result)
	})

	t.Run("Get non-existing boolean environment variable", func(t *testing.T) {
		os.Unsetenv("NON_EXISTING_BOOL")

		result := getEnvBool("NON_EXISTING_BOOL", false)
		assert.False(t, result)
	})
}
