package session

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/logger"
	"github.com/corbado/corbado-go/pkg/validationerror"

	"github.com/corbado/corbado-go/internal/assert"
	entities2 "github.com/corbado/corbado-go/pkg/entities"
	"github.com/corbado/corbado-go/pkg/generated/api"
)

type Session interface {
	ValidateToken(shortSession string) (*entities2.User, error)
}

type Impl struct {
	Client *api.ClientWithResponses
	Config *Config
	Jwks   *keyfunc.JWKS
}

var _ Session = &Impl{}

// New returns new session instance
func New(client *api.ClientWithResponses, config *Config) (*Impl, error) {
	if err := assert.NotNil(client, config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &Impl{
		Client: client,
		Config: config,
	}, nil
}

func newJWKS(config *Config) (*keyfunc.JWKS, error) {
	options := keyfunc.Options{
		RequestFactory: func(_ context.Context, urlAddress string) (*http.Request, error) {
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
		ResponseExtractor: func(_ context.Context, resp *http.Response) (json.RawMessage, error) {
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

func (i *Impl) ValidateToken(shortSession string) (*entities2.User, error) {
	if err := assert.StringNotEmpty(shortSession); err != nil {
		return nil, err
	}

	if i.Jwks == nil {
		jwks, err := newJWKS(i.Config)
		if err != nil {
			return nil, err
		}

		i.Jwks = jwks
	}

	token, err := jwt.ParseWithClaims(shortSession, &entities2.Claims{}, i.Jwks.Keyfunc)
	if err != nil {
		code := validationerror.CodeJWTGeneral
		libraryValidationErr := &jwt.ValidationError{}

		if errors.As(err, &libraryValidationErr) {
			switch {
			case libraryValidationErr.Errors&jwt.ValidationErrorMalformed != 0:
				code = validationerror.CodeJWTInvalidData

			case libraryValidationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				code = validationerror.CodeJWTInvalidSignature

			case libraryValidationErr.Errors&jwt.ValidationErrorNotValidYet != 0:
				code = validationerror.CodeJWTBefore

			case libraryValidationErr.Errors&jwt.ValidationErrorExpired != 0:
				code = validationerror.CodeJWTExpired
			}
		}

		return nil, validationerror.New(err.Error(), code)
	}

	claims := token.Claims.(*entities2.Claims)
	if claims.Issuer != i.Config.JWTIssuer {
		return nil, validationerror.New(
			fmt.Sprintf("JWT issuer mismatch (configured: '%s', actual JWT: '%s')", i.Config.JWTIssuer, claims.Issuer),
			validationerror.CodeJWTIssuerMismatch,
		)
	}

	return &entities2.User{
		UserID:   claims.Subject,
		FullName: claims.Name,
	}, nil
}
