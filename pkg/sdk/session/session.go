package session

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/corbado/corbado-go/pkg/sdk/assert"
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

	projectID              string
	frontendAPI            string
	shortSessionCookieName string
	issuer                 string
}

var _ Session = &Impl{}

// New returns new user client
func New(client *api.ClientWithResponses, projectID string, frontendAPI string, shortSessionCookieName string, issuer string) (*Impl, error) {
	if err := assert.NotNil(client); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(projectID); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(frontendAPI); err != nil {
		return nil, err
	}

	if err := assert.StringNotEmpty(shortSessionCookieName); err != nil {
		return nil, err
	}

	return &Impl{
		client:                 client,
		projectID:              projectID,
		frontendAPI:            frontendAPI,
		shortSessionCookieName: shortSessionCookieName,
		issuer:                 issuer,
	}, nil
}

func (i *Impl) ValidateShortSessionValue(shortSession string) (*entity.User, error) {
	if shortSession == "" {
		return nil, nil
	}

	jwks, err := keyfunc.Get(fmt.Sprintf("%s/.well-known/jwks", i.frontendAPI), keyfunc.Options{
		RequestFactory: func(ctx context.Context, urlAddress string) (*http.Request, error) {
			address, err := url.Parse(urlAddress)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			req := &http.Request{
				Method: http.MethodGet,
				URL:    address,
				Header: map[string][]string{
					"X-Corbado-ProjectID": {i.projectID},
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
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	token, err := jwt.ParseWithClaims(shortSession, &entity.Claims{}, jwks.Keyfunc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claims := token.Claims.(*entity.Claims)
	if i.issuer != "" && claims.Issuer != i.issuer {
		return nil, errors.Errorf("JWT issuer mismatch (configured for Frontend API: '%s', actual JWT: '%s')", i.issuer, claims.Issuer)
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
