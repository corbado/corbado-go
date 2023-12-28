package corbado

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
)

type Configuration struct {
	ProjectID              string
	APISecret              string
	FrontendAPI            string
	BackendAPI             string
	ShortSessionCookieName string
	CacheMaxAge            time.Duration
	JWTIssuer              string

	JWKSRefreshInterval  time.Duration
	JWKSRefreshRateLimit time.Duration
	JWKSRefreshTimeout   time.Duration

	HTTPClient         *http.Client
	ExtraClientOptions []api.ClientOption
}

const (
	configurationDefaultBackendAPI             string = "https://backendapi.corbado.io"
	configurationDefaultFrontendAPI            string = "https://%s.frontendapi.corbado.io"
	configurationDefaultShortSessionCookieName string = "cbo_short_session"
	configurationDefaultCacheMaxAge                   = time.Minute

	configurationDefaultJWKSRefreshInterval  = time.Hour
	configurationDefaultJWKSRefreshRateLimit = 5 * time.Minute
	configurationDefaultJWKSRefreshTimeout   = 10 * time.Second
)

// NewConfiguration returns new configuration with sane defaults
func NewConfiguration(projectID string, apiSecret string) (*Configuration, error) {
	if err := assert.StringNotEmpty(projectID); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(apiSecret); err != nil {
		return nil, err
	}

	return &Configuration{
		ProjectID:              projectID,
		APISecret:              apiSecret,
		BackendAPI:             configurationDefaultBackendAPI,
		FrontendAPI:            fmt.Sprintf(configurationDefaultFrontendAPI, projectID),
		ShortSessionCookieName: configurationDefaultShortSessionCookieName,
		CacheMaxAge:            configurationDefaultCacheMaxAge,
		JWKSRefreshInterval:    configurationDefaultJWKSRefreshInterval,
		JWKSRefreshRateLimit:   configurationDefaultJWKSRefreshRateLimit,
		JWKSRefreshTimeout:     configurationDefaultJWKSRefreshTimeout,
	}, nil
}

// MustNewConfiguration returns new configuration and panics if projectID or apiSecret are not specified/empty
func MustNewConfiguration(projectID string, apiSecret string) *Configuration {
	config, err := NewConfiguration(projectID, apiSecret)
	if err != nil {
		panic(err)
	}

	return config
}
