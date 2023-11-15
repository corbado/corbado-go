package corbado

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
)

type Config struct {
	ProjectID              string
	APISecret              string
	FrontendAPI            string
	BackendAPI             string
	ShortSessionCookieName string
	CacheMaxAge            time.Duration
	JWTIssuer              string

	HTTPClient         *http.Client
	ExtraClientOptions []api.ClientOption
}

const (
	configDefaultBackendAPI             string = "https://backendapi.corbado.io"
	configDefaultFrontendAPI            string = "https://%s.frontendapi.corbado.io"
	configDefaultShortSessionCookieName string = "cbo_short_session"
	configDefaultCacheMaxAge                   = time.Minute
)

// NewConfig returns new config with sane defaults
func NewConfig(projectID string, apiSecret string) (*Config, error) {
	if err := assert.StringNotEmpty(projectID); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(apiSecret); err != nil {
		return nil, err
	}

	return &Config{
		ProjectID:              projectID,
		APISecret:              apiSecret,
		BackendAPI:             configDefaultBackendAPI,
		FrontendAPI:            fmt.Sprintf(configDefaultFrontendAPI, projectID),
		ShortSessionCookieName: configDefaultShortSessionCookieName,
		CacheMaxAge:            configDefaultCacheMaxAge,
	}, nil
}

// MustNewConfig returns new config and panics if projectID or apiSecret are not specified/empty
func MustNewConfig(projectID string, apiSecret string) *Config {
	config, err := NewConfig(projectID, apiSecret)
	if err != nil {
		panic(err)
	}

	return config
}
