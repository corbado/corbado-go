package session

import (
	"time"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/assert"
)

type Config struct {
	ProjectID            string
	JWTIssuer            string
	JwksURI              string
	JWKSRefreshInterval  time.Duration
	JWKSRefreshRateLimit time.Duration
	JWKSRefreshTimeout   time.Duration
}

func (c *Config) validate() error {
	if err := assert.ValidProjectID(c.ProjectID); err != nil {
		return errors.WithMessage(err, "Invalid ProjectID given")
	}

	if err := assert.ValidAPIEndpoint(c.JWTIssuer); err != nil {
		return errors.WithMessage(err, "Invalid JWTIssuer given")
	}

	if err := assert.StringNotEmpty(c.JwksURI); err != nil {
		return errors.WithMessage(err, "Invalid JwksURI given")
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
