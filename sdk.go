package corbado

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/v2/internal/assert"
	"github.com/corbado/corbado-go/v2/internal/services/identifier"
	"github.com/corbado/corbado-go/v2/internal/services/session"
	"github.com/corbado/corbado-go/v2/internal/services/user"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
	"github.com/corbado/corbado-go/v2/pkg/servererror"
	"github.com/corbado/corbado-go/v2/pkg/validationerror"

	// This blank import keeps github.com/go-git/go-git/v5 in go.mod
	// so Dependabot can see and alert on it, but it won't be built
	// into normal binaries.
	_ "github.com/go-git/go-git/v5"
)

const Version = "2.2.2"

type SDK interface {
	Sessions() session.Session
	Users() user.User
	Identifiers() identifier.Identifier
}

type Impl struct {
	client     *api.ClientWithResponses
	HTTPClient *http.Client

	sessions    session.Session
	users       user.User
	identifiers identifier.Identifier
}

var _ SDK = &Impl{}

// NewSDK returns new SDK
func NewSDK(config *Config) (*Impl, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, err
	}

	// instantiate all APIs eagerly because it's cheap to do so and we don't have to deal with thread safety this way

	sessionConfig := &session.Config{
		ProjectID:            config.ProjectID,
		JWTIssuer:            config.FrontendAPI,
		JwksURI:              fmt.Sprintf("%s/.well-known/jwks", config.FrontendAPI),
		JWKSRefreshInterval:  config.JWKSRefreshInterval,
		JWKSRefreshRateLimit: config.JWKSRefreshRateLimit,
		JWKSRefreshTimeout:   config.JWKSRefreshTimeout,
	}

	sessions, err := session.New(client, sessionConfig)
	if err != nil {
		return nil, err
	}

	users, err := user.New(client)
	if err != nil {
		return nil, err
	}

	identifiers, err := identifier.New(client)
	if err != nil {
		return nil, err
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Impl{
		client:      client,
		sessions:    sessions,
		users:       users,
		HTTPClient:  httpClient,
		identifiers: identifiers,
	}, nil
}

// Sessions returns sessions client
func (i *Impl) Sessions() session.Session {
	return i.sessions
}

// Users returns users client
func (i *Impl) Users() user.User {
	return i.users
}

// Identifiers returns identifiers client
func (i *Impl) Identifiers() identifier.Identifier {
	return i.identifiers
}

// IsServerError checks if given error is a ServerError
func IsServerError(err error) bool {
	var serverErr *servererror.ServerError

	return errors.As(err, &serverErr)
}

// AsServerError casts given error into a ServerError, if possible
func AsServerError(err error) *servererror.ServerError {
	var serverErr *servererror.ServerError
	ok := errors.As(err, &serverErr)
	if !ok {
		return nil
	}

	return serverErr
}

// IsValidationError checks if given error is a ValidationError
func IsValidationError(err error) bool {
	var validationErr *validationerror.ValidationError

	return errors.As(err, &validationErr)
}

// AsValidationError casts given error into a ValidationError, if possible
func AsValidationError(err error) *validationerror.ValidationError {
	var validationErr *validationerror.ValidationError
	ok := errors.As(err, &validationErr)
	if !ok {
		return nil
	}

	return validationErr
}
