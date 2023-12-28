package corbado

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/pkg/assert"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/logger"
)

func newClient(config *Configuration) (*api.ClientWithResponses, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth(config.ProjectID, config.APISecret)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	extraOptions := []api.ClientOption{
		api.WithRequestEditorFn(newSDKVersionHeaderEditorFn),
		api.WithRequestEditorFn(basicAuth.Intercept),
	}

	if config.ExtraClientOptions != nil {
		extraOptions = append(extraOptions, config.ExtraClientOptions...)
	}

	return api.NewClientWithResponses(config.BackendAPI, extraOptions...)
}

type httpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type loggingClient struct {
	underlying httpRequestDoer
}

// Do implements HttpRequestDoer and executes HTTP request
func (l *loggingClient) Do(req *http.Request) (*http.Response, error) {
	if err := assert.NotNil(req); err != nil {
		return nil, err
	}

	logger.Debug("Sending request to Public API: %s %s", req.Method, req.URL.String())
	if req.Body != nil {
		requestBody, err := l.readBody(&req.Body)
		if err != nil {
			return nil, err
		}

		logger.Debug("Request body: %s", requestBody)
	}

	response, err := l.underlying.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := l.readBody(&response.Body)
	if err != nil {
		return nil, err
	}

	logger.Debug("Received response from Public API: %s", responseBody)

	return response, nil
}

func (l *loggingClient) readBody(rc *io.ReadCloser /* nilable */) (string, error) {
	if rc == nil {
		return "", nil
	}

	var buf bytes.Buffer
	tee := io.TeeReader(*rc, &buf)

	body, err := io.ReadAll(tee)
	if err != nil {
		return "", errors.WithStack(err)
	}

	*rc = io.NopCloser(&buf)

	return string(body), nil
}

// newLoggingClient returns new logging HTTP client
func newLoggingClient() (*loggingClient, error) {
	return &loggingClient{&http.Client{}}, nil
}

// NewLoggingClientOption enhances HTTP client to log requests/responses
func NewLoggingClientOption() api.ClientOption {
	return func(c *api.Client) error {
		client, err := newLoggingClient()
		if err != nil {
			return err
		}

		c.Client = client

		return nil
	}
}

func newSDKVersionHeaderEditorFn(_ context.Context, req *http.Request) error {
	req.Header.Set("X-Corbado-SDK-Version", fmt.Sprintf("Go SDK %s", Version))

	return nil
}
