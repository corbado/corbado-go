package corbado

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"
	"github.com/corbado/corbado-go/internal/services/identifier"
	"github.com/corbado/corbado-go/internal/services/session"
	"github.com/corbado/corbado-go/internal/services/user"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

const Version = "2.0.0"

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

// Users returns identifiers client
func (i *Impl) Identifiers() identifier.Identifier {
	return i.Identifiers()
}

// IsServerError checks if given error is a ServerError
func IsServerError(err error) bool {
	var serverError *servererror.ServerError
	ok := errors.As(err, &serverError)

	return ok
}

// AsServerError casts given error into a ServerError, if possible
func AsServerError(err error) *servererror.ServerError {
	var serverError *servererror.ServerError
	ok := errors.As(err, &serverError)
	if !ok {
		return nil
	}

	return serverError
}
