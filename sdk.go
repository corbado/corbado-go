package corbado

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/pkg/services/authtoken"
	"github.com/corbado/corbado-go/pkg/services/emailmagiclink"
	"github.com/corbado/corbado-go/pkg/services/emailotp"
	"github.com/corbado/corbado-go/pkg/services/passkey"
	"github.com/corbado/corbado-go/pkg/services/project"
	"github.com/corbado/corbado-go/pkg/services/session"
	"github.com/corbado/corbado-go/pkg/services/smsotp"
	"github.com/corbado/corbado-go/pkg/services/template"
	"github.com/corbado/corbado-go/pkg/services/user"
	"github.com/corbado/corbado-go/pkg/services/validation"
)

const Version = "1.0.0"

type SDK interface {
	AuthTokens() authtoken.AuthToken
	EmailMagicLinks() emailmagiclink.EmailMagicLink
	EmailOTPs() emailotp.EmailOTP
	Passkeys() passkey.Passkey
	Projects() project.Project
	Sessions() session.Session
	SmsOTPs() smsotp.SmsOTP
	Templates() template.Template
	Users() user.User
	Validations() validation.Validation
}

type Impl struct {
	client     *api.ClientWithResponses
	HTTPClient *http.Client

	authTokens      authtoken.AuthToken
	emailMagicLinks emailmagiclink.EmailMagicLink
	emailOTPs       emailotp.EmailOTP
	passkeys        passkey.Passkey
	projects        project.Project
	sessions        session.Session
	smsOTPs         smsotp.SmsOTP
	templates       template.Template
	users           user.User
	validations     validation.Validation
}

var _ SDK = &Impl{}

// NewSDK returns new SDK
func NewSDK(config *Config) (*Impl, error) {
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

	emailMagicLinks, err := emailmagiclink.New(client)
	if err != nil {
		return nil, err
	}

	emailOTPs, err := emailotp.New(client)
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

	sessionConfig := &session.Config{
		ProjectID:            config.ProjectID,
		FrontendAPI:          config.FrontendAPI,
		JWTIssuer:            config.JWTIssuer,
		JWKSRefreshInterval:  config.JWKSRefreshInterval,
		JWKSRefreshRateLimit: config.JWKSRefreshRateLimit,
		JWKSRefreshTimeout:   config.JWKSRefreshTimeout,
	}

	sessions, err := session.New(client, sessionConfig)
	if err != nil {
		return nil, err
	}

	smsOTPs, err := smsotp.New(client)
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

	validations, err := validation.New(client)
	if err != nil {
		return nil, err
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Impl{
		client:          client,
		authTokens:      authTokens,
		emailMagicLinks: emailMagicLinks,
		emailOTPs:       emailOTPs,
		passkeys:        passkeys,
		projects:        projects,
		sessions:        sessions,
		smsOTPs:         smsOTPs,
		templates:       templates,
		users:           users,
		validations:     validations,
		HTTPClient:      httpClient,
	}, nil
}

// AuthTokens returns auth tokens client
func (i *Impl) AuthTokens() authtoken.AuthToken {
	return i.authTokens
}

// EmailMagicLinks returns email magic links client
func (i *Impl) EmailMagicLinks() emailmagiclink.EmailMagicLink {
	return i.emailMagicLinks
}

// EmailOTPs returns email OTPs client
func (i *Impl) EmailOTPs() emailotp.EmailOTP {
	return i.emailOTPs
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

// SmsOTPs returns sms OTPs client
func (i *Impl) SmsOTPs() smsotp.SmsOTP {
	return i.smsOTPs
}

// Templates returns templates client
func (i *Impl) Templates() template.Template {
	return i.templates
}

// Users returns users client
func (i *Impl) Users() user.User {
	return i.users
}

// Validations returns validation client
func (i *Impl) Validations() validation.Validation {
	return i.validations
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
