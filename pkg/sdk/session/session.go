package session

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/corbado/corbado-go/pkg/logger"
	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/config"
	"github.com/corbado/corbado-go/pkg/sdk/entity"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Session interface {
	ValidateShortSessionValue(shortSession string) (*entity.User, error)
	GetCurrentUser(shortSession string) (*entity.User, error)
}

type Impl struct {
	client *api.ClientWithResponses
	config *config.Config
	jwks   *keyfunc.JWKS
}

var _ Session = &Impl{}

// New returns new user client
func New(client *api.ClientWithResponses, config *config.Config) (*Impl, error) {
	if err := assert.NotNil(client, config); err != nil {
		return nil, err
	}

	jwks, err := newJWKS(config)
	if err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
		config: config,
		jwks:   jwks,
	}, nil
}

func newJWKS(config *config.Config) (*keyfunc.JWKS, error) {
	options := keyfunc.Options{
		RequestFactory: func(ctx context.Context, urlAddress string) (*http.Request, error) {
			address, err := url.Parse(urlAddress)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			req := &http.Request{
				Method: http.MethodGet,
				URL:    address,
				Header: map[string][]string{
					"X-Corbado-ProjectID": {config.ProjectID},
				},
			}

			return req, nil
		},
		ResponseExtractor: func(ctx context.Context, resp *http.Response) (json.RawMessage, error) {
			rspBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return rspBody, nil
		},
		RefreshErrorHandler: func(err error) {
			logger.Error("Error refreshing JWKS: %s", err.Error())
		},
		RefreshInterval:   config.JWKSRefreshInterval,
		RefreshRateLimit:  config.JWKSRefreshRateLimit,
		RefreshTimeout:    config.JWKSRefreshTimeout,
		RefreshUnknownKID: true,
	}

	return keyfunc.Get(fmt.Sprintf("%s/.well-known/jwks", config.FrontendAPI), options)
}

func (i *Impl) ValidateShortSessionValue(shortSession string) (*entity.User, error) {
	if shortSession == "" {
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(shortSession, &entity.Claims{}, i.jwks.Keyfunc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claims := token.Claims.(*entity.Claims)
	if i.config.JWTIssuer != "" && claims.Issuer != i.config.JWTIssuer {
		return nil, errors.Errorf("JWT issuer mismatch (configured for Frontend API: '%s', actual JWT: '%s')", i.config.JWTIssuer, claims.Issuer)
	}

	return &entity.User{
		Authenticated: true,
		ID:            claims.Subject,
		Name:          claims.Name,
		Email:         claims.Email,
		PhoneNumber:   claims.PhoneNumber,
	}, nil
}

func (i *Impl) GetCurrentUser(shortSession string) (*entity.User, error) {
	usr, err := i.ValidateShortSessionValue(shortSession)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		return usr, nil
	}

	return entity.NewGuestUser(), nil
}
