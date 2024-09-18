package corbado

import (
	"testing"

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
