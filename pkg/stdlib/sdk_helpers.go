package stdlib

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/internal/assert"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/common"
	"github.com/corbado/corbado-go/pkg/util"
)

type SDKHelpers struct {
	config *corbado.Config
}

func NewSDKHelpers(config *corbado.Config) (*SDKHelpers, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	return &SDKHelpers{
		config: config,
	}, nil
}

func (s *SDKHelpers) GetShortSessionValue(req *http.Request, shortSessionCookieName string) (string, error) {
	if err := assert.NotNil(req); err != nil {
		return "", err
	}

	if err := assert.StringNotEmpty(shortSessionCookieName); err != nil {
		return "", err
	}

	ses, err := s.getShortSessionValueFromCookie(req, shortSessionCookieName)
	if err != nil {
		return "", err
	}

	if ses == "" {
		return s.getShortSessionValueFromAuthHeader(req)
	}

	return ses, nil
}

func (s *SDKHelpers) getShortSessionValueFromCookie(req *http.Request, shortSessionCookieName string) (string, error) {
	if err := assert.NotNil(req); err != nil {
		return "", err
	}

	cookie, err := req.Cookie(shortSessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (s *SDKHelpers) getShortSessionValueFromAuthHeader(req *http.Request) (string, error) {
	if err := assert.NotNil(req); err != nil {
		return "", err
	}

	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 8 {
		return "", nil
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", nil
	}

	return authHeader[7:], nil
}

func (s *SDKHelpers) GetClientInfo(req *http.Request) (*common.ClientInfo, error) {
	if err := assert.NotNil(req); err != nil {
		return nil, err
	}

	ip, err := s.GetRemoteAddress(req)
	if err != nil {
		return nil, err
	}

	return util.ClientInfo(req.UserAgent(), ip), nil
}

func (s *SDKHelpers) GetRemoteAddress(req *http.Request) (string, error) {
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip, nil
	}

	return req.RemoteAddr, nil
}
