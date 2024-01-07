package session

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"
	"github.com/corbado/corbado-go/internal/entities"
	"github.com/corbado/corbado-go/internal/logger"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
)

type Session interface {
	ValidateShortSessionValue(shortSession string) (*entities.User, error)
	GetCurrentUser(shortSession string) (*entities.User, error)
	ConfigGet(ctx context.Context, params *api.SessionConfigGetParams, editors ...api.RequestEditorFn) (*api.SessionConfigGetRsp, error)
	LongSessionRevoke(ctx context.Context, sessionID string, req api.LongSessionRevokeReq, editors ...api.RequestEditorFn) error
	LongSessionGet(ctx context.Context, sessionID string, editors ...api.RequestEditorFn) (*api.LongSessionGetRsp, error)
}

type Impl struct {
	client *api.ClientWithResponses
	config *Config
	jwks   *keyfunc.JWKS
}

var _ Session = &Impl{}

// New returns new user client
func New(client *api.ClientWithResponses, config *Config) (*Impl, error) {
	if err := assert.NotNil(client, config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &Impl{
		client: client,
		config: config,
	}, nil
}

func newJWKS(config *Config) (*keyfunc.JWKS, error) {
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

	return keyfunc.Get(config.JwksURI, options)
}

func (i *Impl) ValidateShortSessionValue(shortSession string) (*entities.User, error) {
	if shortSession == "" {
		return nil, nil
	}

	if i.jwks == nil {
		jwks, err := newJWKS(i.config)
		if err != nil {
			return nil, err
		}

		i.jwks = jwks
	}

	token, err := jwt.ParseWithClaims(shortSession, &entities.Claims{}, i.jwks.Keyfunc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claims := token.Claims.(*entities.Claims)
	if claims.Issuer != i.config.JWTIssuer {
		return nil, errors.Errorf("JWT issuer mismatch (configured: '%s', actual JWT: '%s')", i.config.JWTIssuer, claims.Issuer)
	}

	return &entities.User{
		Authenticated: true,
		ID:            claims.Subject,
		Name:          claims.Name,
		Email:         claims.Email,
		PhoneNumber:   claims.PhoneNumber,
	}, nil
}

func (i *Impl) GetCurrentUser(shortSession string) (*entities.User, error) {
	usr, err := i.ValidateShortSessionValue(shortSession)
	if err != nil {
		return nil, err
	}

	if usr != nil {
		return usr, nil
	}

	return entities.NewGuestUser(), nil
}

// ConfigGet retrieves session config by projectID inferred from authentication
func (i *Impl) ConfigGet(ctx context.Context, params *api.SessionConfigGetParams, editors ...api.RequestEditorFn) (*api.SessionConfigGetRsp, error) {
	res, err := i.client.SessionConfigGetWithResponse(ctx, params, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}

// LongSessionRevoke revokes an active long session by sessionID
func (i *Impl) LongSessionRevoke(ctx context.Context, sessionID string, req api.LongSessionRevokeReq, editors ...api.RequestEditorFn) error {
	res, err := i.client.LongSessionRevokeWithResponse(ctx, sessionID, req, editors...)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return servererror.New(res.JSONDefault)
	}

	return nil
}

// LongSessionGet gets a long session by sessionID
func (i *Impl) LongSessionGet(ctx context.Context, sessionID string, editors ...api.RequestEditorFn) (*api.LongSessionGetRsp, error) {
	res, err := i.client.LongSessionGetWithResponse(ctx, sessionID, editors...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.JSONDefault != nil {
		return nil, servererror.New(res.JSONDefault)
	}

	return res.JSON200, nil
}
