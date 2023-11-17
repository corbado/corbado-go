package stdlib

import (
	"net/http"
	"strings"

	"github.com/corbado/corbado-go/pkg/sdk/assert"
	"github.com/corbado/corbado-go/pkg/sdk/config"
	"github.com/corbado/corbado-go/pkg/sdk/entity/common"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/pkg/errors"
)

type SDKHelpers struct {
	config *config.Config
}

func NewSDKHelpers(config *config.Config) (*SDKHelpers, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	return &SDKHelpers{
		config: config,
	}, nil
}

func (s *SDKHelpers) GetShortSessionValue(req *http.Request) (string, error) {
	if err := assert.NotNil(req); err != nil {
		return "", err
	}

	ses, err := s.getShortSessionValueFromCookie(req)
	if err != nil {
		return "", err
	}

	if ses == "" {
		return s.getShortSessionValueFromAuthHeader(req)
	}

	return "", nil
}

func (s *SDKHelpers) getShortSessionValueFromCookie(req *http.Request) (string, error) {
	if err := assert.NotNil(req); err != nil {
		return "", err
	}

	cookie, err := req.Cookie(s.config.ShortSessionCookieName)
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

	return util.ClientInfo(ip, req.UserAgent()), nil
}

func (s *SDKHelpers) GetRemoteAddress(req *http.Request) (string, error) {
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip, nil
	}

	return req.RemoteAddr, nil
}
