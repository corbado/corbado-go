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

	"github.com/corbado/corbado-go/v2/pkg/logger"
	"github.com/corbado/corbado-go/v2/pkg/validationerror"

	"github.com/corbado/corbado-go/v2/internal/assert"
	"github.com/corbado/corbado-go/v2/pkg/entities"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
)

type Session interface {
	ValidateToken(sessionToken string) (*entities.User, error)
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

func (i *Impl) ValidateToken(sessionToken string) (*entities.User, error) {
	if err := assert.StringNotEmpty(sessionToken); err != nil {
		return nil, err
	}

	if i.Jwks == nil {
		jwks, err := newJWKS(i.Config)
		if err != nil {
			return nil, err
		}

		i.Jwks = jwks
	}

	token, err := jwt.ParseWithClaims(sessionToken, &entities.Claims{}, i.Jwks.Keyfunc)
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

		return nil, newValidationError(err.Error(), sessionToken, code)
	}

	claims := token.Claims.(*entities.Claims)
	if err := i.validateIssuer(claims.Issuer, sessionToken); err != nil {
		return nil, err
	}

	return &entities.User{
		UserID:   claims.Subject,
		FullName: claims.Name,
	}, nil
}

func (i *Impl) validateIssuer(jwtIssuer string, sessionToken string) error {
	if jwtIssuer == "" {
		return newValidationError("Issuer is empty", sessionToken, validationerror.CodeJWTIssuerEmpty)
	}

	// Compare to old Frontend API (without .cloud.) to make our Frontend API host name change downwards compatible
	if jwtIssuer == fmt.Sprintf("https://%s.frontendapi.corbado.io", i.Config.ProjectID) {
		return nil
	}

	// Compare to new Frontend API (with .cloud.)
	if jwtIssuer == fmt.Sprintf("https://%s.frontendapi.cloud.corbado.io", i.Config.ProjectID) {
		return nil
	}

	// Compare to configured issuer (from FrontendAPI), needed if you set a CNAME for example
	if jwtIssuer != i.Config.JWTIssuer {
		return newValidationError(
			fmt.Sprintf("Issuer mismatch (configured trough FrontendAPI: '%s', JWT issuer: '%s')", i.Config.JWTIssuer, jwtIssuer),
			sessionToken,
			validationerror.CodeJWTIssuerMismatch,
		)
	}

	return nil
}

func newValidationError(message string, jwt string, code validationerror.Code) error {
	return validationerror.New(fmt.Sprintf("JWT validation failed: '%s' (JWT: '%s')", message, jwt), code)
}
