package corbado

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
)

type Config struct {
	ProjectID              string
	APISecret              string
	FrontendAPI            string
	BackendAPI             string
	ShortSessionCookieName string
	CacheMaxAge            time.Duration

	JWKSRefreshInterval  time.Duration
	JWKSRefreshRateLimit time.Duration
	JWKSRefreshTimeout   time.Duration

	HTTPClient         *http.Client
	ExtraClientOptions []api.ClientOption
}

const (
	configDefaultBackendAPI             string = "https://backendapi.corbado.io"
	configDefaultFrontendAPI            string = "https://%s.frontendapi.corbado.io"
	configDefaultShortSessionCookieName string = "cbo_short_session"
	configDefaultCacheMaxAge                   = time.Minute

	configDefaultJWKSRefreshInterval  = time.Hour
	configDefaultJWKSRefreshRateLimit = 5 * time.Minute
	configDefaultJWKSRefreshTimeout   = 10 * time.Second
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
		FrontendAPI:            fmt.Sprintf(configDefaultFrontendAPI, projectID),
		BackendAPI:             configDefaultBackendAPI,
		ShortSessionCookieName: configDefaultShortSessionCookieName,
		CacheMaxAge:            configDefaultCacheMaxAge,
		JWKSRefreshInterval:    configDefaultJWKSRefreshInterval,
		JWKSRefreshRateLimit:   configDefaultJWKSRefreshRateLimit,
		JWKSRefreshTimeout:     configDefaultJWKSRefreshTimeout,
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

// NewConfigEnv returns new config with values from env variables (CORBADO_PROJECT_ID and CORBADO_API_SECRET)
func NewConfigEnv() (*Config, error) {
	projectID := os.Getenv("CORBADO_PROJECT_ID")
	if projectID == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_PROJECT_ID")
	}

	apiSecret := os.Getenv("CORBADO_API_SECRET")
	if apiSecret == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_API_SECRET")
	}

	return NewConfig(projectID, apiSecret)
}

func (c *Config) validate() error {
	if err := assert.ValidProjectID(c.ProjectID); err != nil {
		return errors.WithMessage(err, "Invalid ProjectID given")
	}

	if err := assert.ValidAPISecret(c.APISecret); err != nil {
		return errors.WithMessage(err, "Invalid APISecret given")
	}

	if err := assert.ValidAPIEndpoint(c.FrontendAPI); err != nil {
		return errors.WithMessage(err, "Invalid FrontendAPI given")
	}

	if err := assert.ValidAPIEndpoint(c.BackendAPI); err != nil {
		return errors.WithMessage(err, "Invalid BackendAPI given")
	}

	if err := assert.StringNotEmpty(c.ShortSessionCookieName); err != nil {
		return errors.WithMessage(err, "Invalid ShortSessionCookieName given")
	}

	if err := assert.DurationNotEmpty(c.CacheMaxAge); err != nil {
		return errors.WithMessage(err, "Invalid CacheMaxAge given")
	}

	if err := assert.DurationNotEmpty(c.JWKSRefreshInterval); err != nil {
		return errors.WithMessage(err, "Invalid JWKSRefreshInterval given")
	}

	if err := assert.DurationNotEmpty(c.JWKSRefreshRateLimit); err != nil {
		return errors.WithMessage(err, "Invalid JWKSRefreshRateLimit given")
	}

	if err := assert.DurationNotEmpty(c.JWKSRefreshTimeout); err != nil {
		return errors.WithMessage(err, "Invalid JWKSRefreshTimeout given")
	}

	return nil
}
