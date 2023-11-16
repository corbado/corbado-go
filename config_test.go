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

	cfg, err := NewConfig(projectID, secret)
	require.NoError(t, err)

	assert.Equal(t, projectID, cfg.ProjectID)
	assert.Equal(t, secret, cfg.APISecret)
	assert.Equal(t, fmt.Sprintf(configDefaultFrontendAPI, projectID), cfg.FrontendAPI)
	assert.Equal(t, configDefaultBackendAPI, cfg.BackendAPI)
	assert.Equal(t, configDefaultShortSessionCookieName, cfg.ShortSessionCookieName)
	assert.Equal(t, configDefaultCacheMaxAge, cfg.CacheMaxAge)
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
			cfg, err := NewConfig(test.projectID, test.secret)
			assert.Nil(t, cfg)
			assert.ErrorContains(t, err, "given value '' is too short")
		})
	}
}
