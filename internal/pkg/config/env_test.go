package config

import (
	"os"
	"testing"

	"github.com/caarlos0/env/v9"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	envVars := map[string]string{
		"APP_SECRET_KEY":         "test-secret",
		"DATABASE_HOST":          "localhost",
		"DATABASE_PORT":          "5432",
		"DATABASE_USER":          "testuser",
		"DATABASE_PASSWORD":      "testpass",
		"DATABASE_NAME":          "testdb",
		"DATABASE_LOG_LEVEL":     "debug",
		"REDIS_URL":              "redis://localhost:6379",
		"QISCUS_APP_ID":          "test-app-id",
		"QISCUS_SECRET_KEY":      "test-qiscus-secret",
		"QISCUS_OMNICHANNEL_URL": "https://test.qiscus.com",
	}

	for k, v := range envVars {
		t.Setenv(k, v)
	}

	config := Load()

	assert.Equal(t, "test-secret", config.App.SecretKey)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.Name)
	assert.Equal(t, "debug", config.Database.LogLevel)
	assert.Equal(t, "redis://localhost:6379", config.Redis.URL)
	assert.Equal(t, "test-app-id", config.Qiscus.AppID)
	assert.Equal(t, "test-qiscus-secret", config.Qiscus.SecretKey)
	assert.Equal(t, "https://test.qiscus.com", config.Qiscus.Omnichannel.URL)
}

func TestDatabase_DataSourceName(t *testing.T) {
	tests := []struct {
		name     string
		database Database
		expected string
	}{
		{
			name: "Complete database config",
			database: Database{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
			},
			expected: "user=testuser password=testpass host=localhost port=5432 dbname=testdb sslmode=disable",
		},
		{
			name: "Empty password",
			database: Database{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "",
				Name:     "testdb",
			},
			expected: "user=testuser password= host=localhost port=5432 dbname=testdb sslmode=disable",
		},
		{
			name: "Different port",
			database: Database{
				Host:     "localhost",
				Port:     5433,
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
			},
			expected: "user=testuser password=testpass host=localhost port=5433 dbname=testdb sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.database.DataSourceName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    map[string]string
		expectError bool
	}{
		{
			name: "Valid configuration",
			setupEnv: map[string]string{
				"APP_SECRET_KEY":         "secret",
				"DATABASE_HOST":          "localhost",
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "user",
				"DATABASE_PASSWORD":      "pass",
				"DATABASE_NAME":          "dbname",
				"DATABASE_LOG_LEVEL":     "",
				"REDIS_URL":              "redis://localhost",
				"QISCUS_APP_ID":          "appid",
				"QISCUS_SECRET_KEY":      "secret",
				"QISCUS_OMNICHANNEL_URL": "https://qiscus.com",
			},
			expectError: false,
		},
		{
			name: "Invalid database port",
			setupEnv: map[string]string{
				"DATABASE_PORT": "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k := range tt.setupEnv {
				os.Unsetenv(k)
			}

			for k, v := range tt.setupEnv {
				t.Setenv(k, v)
			}

			var c Config
			err := env.Parse(&c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
