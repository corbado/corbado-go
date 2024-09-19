package corbado

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Success(t *testing.T) {
	projectID := "pro-12345678"
	secret := "mysupersecret"
	backendAPI := "http://localhost:8080"
	frontendAPI := "http://localhost:8081"

	cfg, err := NewConfig(projectID, secret, frontendAPI, backendAPI)
	require.NoError(t, err)

	assert.Equal(t, projectID, cfg.ProjectID)
	assert.Equal(t, secret, cfg.APISecret)
	assert.Equal(t, frontendAPI, cfg.FrontendAPI)
	assert.Equal(t, backendAPI, cfg.BackendAPI)
	assert.Equal(t, configDefaultCacheMaxAge, cfg.CacheMaxAge)
	assert.Equal(t, configDefaultJWKSRefreshInterval, cfg.JWKSRefreshInterval)
	assert.Equal(t, configDefaultJWKSRefreshRateLimit, cfg.JWKSRefreshRateLimit)
	assert.Equal(t, configDefaultJWKSRefreshTimeout, cfg.JWKSRefreshTimeout)
}

func TestNewConfig_Failure(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		secret      string
		frontendAPI string
		backendAPI  string
	}{
		{
			name: "empty projectID and secret",
		},
		{
			name:   "empty projectID",
			secret: "secret",
		},
		{
			name:      "empty secret",
			projectID: "pro-12345678",
		},
		{
			name:      "empty frontendAPI",
			projectID: "pro-12345678",
			secret:    "secret",
		},
		{
			name:        "empty backendAPI",
			projectID:   "pro-12345678",
			secret:      "secret",
			frontendAPI: "http://localhost:8080",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := NewConfig(test.projectID, test.secret, test.frontendAPI, test.backendAPI)
			assert.Nil(t, cfg)
			assert.ErrorContains(t, err, "given value '' is too short")
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name                  string
		config                Config
		expectedErrorContains string
	}{
		{
			name:                  "invalid ProjectID",
			config:                Config{ProjectID: "", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090"},
			expectedErrorContains: "Invalid ProjectID given",
		},
		{
			name:                  "invalid APISecret",
			config:                Config{ProjectID: "pro-12345678", APISecret: "", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090"},
			expectedErrorContains: "Invalid APISecret given",
		},
		{
			name:                  "invalid FrontendAPI",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "", BackendAPI: "http://localhost:9090"},
			expectedErrorContains: "Invalid FrontendAPI given",
		},
		{
			name:                  "invalid BackendAPI",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: ""},
			expectedErrorContains: "Invalid BackendAPI given",
		},
		{
			name:                  "invalid CacheMaxAge",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090", CacheMaxAge: 0},
			expectedErrorContains: "Invalid CacheMaxAge given",
		},
		{
			name:                  "invalid JWKSRefreshInterval",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090", CacheMaxAge: 10 * time.Second, JWKSRefreshInterval: 0},
			expectedErrorContains: "Invalid JWKSRefreshInterval given",
		},
		{
			name:                  "invalid JWKSRefreshRateLimit",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090", CacheMaxAge: 10 * time.Second, JWKSRefreshInterval: 10 * time.Second, JWKSRefreshRateLimit: 0},
			expectedErrorContains: "Invalid JWKSRefreshRateLimit given",
		},
		{
			name:                  "invalid JWKSRefreshTimeout",
			config:                Config{ProjectID: "pro-12345678", APISecret: "corbado1_secret", FrontendAPI: "http://localhost:8080", BackendAPI: "http://localhost:9090", CacheMaxAge: 10 * time.Second, JWKSRefreshInterval: 10 * time.Second, JWKSRefreshRateLimit: 10 * time.Second, JWKSRefreshTimeout: 0},
			expectedErrorContains: "Invalid JWKSRefreshTimeout given",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.validate()

			// Assert that an error exists and contains the expected substring
			if err == nil {
				t.Errorf("Expected an error, but got none")
			} else {
				if !strings.Contains(err.Error(), test.expectedErrorContains) {
					t.Errorf("Expected error containing %q, but got %q", test.expectedErrorContains, err.Error())
				}
			}
		})
	}
}

func TestConfig_Validate_Success(t *testing.T) {
	validConfig := Config{
		ProjectID:            "pro-12345678",
		APISecret:            "corbado1_secret",
		FrontendAPI:          "http://localhost:8080",
		BackendAPI:           "http://localhost:9090",
		CacheMaxAge:          10 * time.Second,
		JWKSRefreshInterval:  10 * time.Second,
		JWKSRefreshRateLimit: 10 * time.Second,
		JWKSRefreshTimeout:   10 * time.Second,
	}

	err := validConfig.validate()

	// Assert no error is returned for valid config
	if err != nil {
		t.Errorf("Expected no error, but got %q", err.Error())
	}
}
