package corbado

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Success(t *testing.T) {
	projectID := "pro-12345678"
	secret := "mysupersecret"

	cfg, err := NewConfiguration(projectID, secret)
	require.NoError(t, err)

	assert.Equal(t, projectID, cfg.ProjectID)
	assert.Equal(t, secret, cfg.APISecret)
	assert.Equal(t, fmt.Sprintf(configurationDefaultFrontendAPI, projectID), cfg.FrontendAPI)
	assert.Equal(t, configurationDefaultBackendAPI, cfg.BackendAPI)
	assert.Equal(t, configurationDefaultShortSessionCookieName, cfg.ShortSessionCookieName)
	assert.Equal(t, configurationDefaultCacheMaxAge, cfg.CacheMaxAge)
	assert.Equal(t, configurationDefaultJWKSRefreshInterval, cfg.JWKSRefreshInterval)
	assert.Equal(t, configurationDefaultJWKSRefreshRateLimit, cfg.JWKSRefreshRateLimit)
	assert.Equal(t, configurationDefaultJWKSRefreshTimeout, cfg.JWKSRefreshTimeout)
}

func TestNewConfig_Failure(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		secret    string
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := NewConfiguration(test.projectID, test.secret)
			assert.Nil(t, cfg)
			assert.ErrorContains(t, err, "given value '' is too short")
		})
	}
}
