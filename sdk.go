package corbado

import (
	"net/http"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/authtoken"
	"github.com/corbado/corbado-go/pkg/sdk/config"
	"github.com/corbado/corbado-go/pkg/sdk/emailcode"
	"github.com/corbado/corbado-go/pkg/sdk/emaillink"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/passkey"
	"github.com/corbado/corbado-go/pkg/sdk/project"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
	"github.com/corbado/corbado-go/pkg/sdk/session"
	"github.com/corbado/corbado-go/pkg/sdk/template"
	"github.com/corbado/corbado-go/pkg/sdk/user"
	"github.com/corbado/corbado-go/pkg/sdk/validation"
	"github.com/pkg/errors"
)

const Version = "v0.6.0"

type SDK interface {
	AuthTokens() authtoken.AuthToken
	EmailCodes() emailcode.EmailCode
	EmailLinks() emaillink.EmailLink
	Passkeys() passkey.Passkey
	Projects() project.Project
	Sessions() session.Session
	Templates() template.Template
	Users() user.User
	Validations() validation.Validation
}

type Impl struct {
	client     *api.ClientWithResponses
	HTTPClient *http.Client

	authTokens authtoken.AuthToken
	emailCodes emailcode.EmailCode
	emailLinks emaillink.EmailLink
	passkeys   passkey.Passkey
	projects   project.Project
	sessions   session.Session
	templates  template.Template
	validation validation.Validation
	users      user.User
}

var _ SDK = &Impl{}

// NewSDK returns new SDK
func NewSDK(config *config.Config) (*Impl, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	client, err := newClient(config)
	if err != nil {
		return nil, err
	}

	// instantiate all APIs eagerly because it's cheap to do so and we don't have to deal with thread safety this way
	authTokens, err := authtoken.New(client)
	if err != nil {
		return nil, err
	}

	emailCodes, err := emailcode.New(client)
	if err != nil {
		return nil, err
	}

	emailLinks, err := emaillink.New(client)
	if err != nil {
		return nil, err
	}

	passkeys, err := passkey.New(client)
	if err != nil {
		return nil, err
	}

	projects, err := project.New(client)
	if err != nil {
		return nil, err
	}

	sessions, err := session.New(client, config)
	if err != nil {
		return nil, err
	}

	templates, err := template.New(client)
	if err != nil {
		return nil, err
	}

	users, err := user.New(client)
	if err != nil {
		return nil, err
	}

	validation, err := validation.New(client)
	if err != nil {
		return nil, err
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Impl{
		client:     client,
		authTokens: authTokens,
		emailCodes: emailCodes,
		emailLinks: emailLinks,
		passkeys:   passkeys,
		projects:   projects,
		sessions:   sessions,
		templates:  templates,
		users:      users,
		validation: validation,
		HTTPClient: httpClient,
	}, nil
}

// AuthTokens returns auth tokens client
func (i *Impl) AuthTokens() authtoken.AuthToken {
	return i.authTokens
}

// EmailCodes returns email codes client
func (i *Impl) EmailCodes() emailcode.EmailCode {
	return i.emailCodes
}

// EmailLinks returns email links client
func (i *Impl) EmailLinks() emaillink.EmailLink {
	return i.emailLinks
}

// Validations returns validation client
func (i *Impl) Validations() validation.Validation {
	return i.validation
}

// Passkeys returns passkeys client
func (i *Impl) Passkeys() passkey.Passkey {
	return i.passkeys
}

// Projects returns projects client
func (i *Impl) Projects() project.Project {
	return i.projects
}

// Sessions returns sessions client
func (i *Impl) Sessions() session.Session {
	return i.sessions
}

// Templates returns templates client
func (i *Impl) Templates() template.Template {
	return i.templates
}

// Users returns users client
func (i *Impl) Users() user.User {
	return i.users
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

// NewConfig returns new config with sane defaults
// this is a convenience function for config.NewConfig to avoid import cycles
func NewConfig(projectID string, apiSecret string) (*config.Config, error) {
	return config.NewConfig(projectID, apiSecret)
}

// MustNewConfig returns new config and panics if projectID or apiSecret are not specified/empty
// this is a convenience function for config.NewConfig to avoid import cycles
func MustNewConfig(projectID string, apiSecret string) *config.Config {
	return config.MustNewConfig(projectID, apiSecret)
}
