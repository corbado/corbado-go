package corbado

import (
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
)

type Config struct {
	ProjectID   string
	APISecret   string
	FrontendAPI string
	BackendAPI  string
	CacheMaxAge time.Duration

	JWKSRefreshInterval  time.Duration
	JWKSRefreshRateLimit time.Duration
	JWKSRefreshTimeout   time.Duration

	HTTPClient         *http.Client
	ExtraClientOptions []api.ClientOption
}

const (
	configDefaultCacheMaxAge = time.Minute

	configDefaultJWKSRefreshInterval  = time.Hour
	configDefaultJWKSRefreshRateLimit = 5 * time.Minute
	configDefaultJWKSRefreshTimeout   = 10 * time.Second
)

// NewConfig returns new config with sane defaults
func NewConfig(projectID string, apiSecret string, frontendApi string, backendAPI string) (*Config, error) {
	if err := assert.StringNotEmpty(projectID); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(apiSecret); err != nil {
		return nil, err
	}

	return &Config{
		ProjectID:            projectID,
		APISecret:            apiSecret,
		FrontendAPI:          frontendApi,
		BackendAPI:           backendAPI,
		CacheMaxAge:          configDefaultCacheMaxAge,
		JWKSRefreshInterval:  configDefaultJWKSRefreshInterval,
		JWKSRefreshRateLimit: configDefaultJWKSRefreshRateLimit,
		JWKSRefreshTimeout:   configDefaultJWKSRefreshTimeout,
	}, nil
}

// MustNewConfig returns new config and panics if projectID, apiSecret, frontendApi or backendApi are not specified/empty
func MustNewConfig(projectID string, apiSecret string, frontendApi string, backendApi string) *Config {
	config, err := NewConfig(projectID, apiSecret, frontendApi, backendApi)
	if err != nil {
		panic(err)
	}

	return config
}

// NewConfigFromEnv returns new config with values from env variables (CORBADO_PROJECT_ID and CORBADO_API_SECRET)
func NewConfigFromEnv() (*Config, error) {
	projectID := os.Getenv("CORBADO_PROJECT_ID")
	if projectID == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_PROJECT_ID")
	}

	apiSecret := os.Getenv("CORBADO_API_SECRET")
	if apiSecret == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_API_SECRET")
	}

	frontendApi := os.Getenv("CORBADO_FRONTEND_API")
	if frontendApi == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_FRONTEND_API")
	}

	backendApi := os.Getenv("CORBADO_BACKEND_API")
	if backendApi == "" {
		return nil, errors.Errorf("Missing env variable CORBADO_BACKEND_API")
	}

	return NewConfig(projectID, apiSecret, frontendApi, backendApi)
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
