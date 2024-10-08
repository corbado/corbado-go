package corbado

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/pkg/errors"

	"github.com/corbado/corbado-go/v2/pkg/logger"

	"github.com/corbado/corbado-go/v2/internal/assert"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
)

func newClient(config *Config) (*api.ClientWithResponses, error) {
	if err := assert.NotNil(config); err != nil {
		return nil, err
	}

	basicAuth, err := securityprovider.NewSecurityProviderBasicAuth(config.ProjectID, config.APISecret)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	extraOptions := []api.ClientOption{
		api.WithRequestEditorFn(newSDKHeaderEditorFn),
		api.WithRequestEditorFn(basicAuth.Intercept),
	}

	if config.ExtraClientOptions != nil {
		extraOptions = append(extraOptions, config.ExtraClientOptions...)
	}

	backendServer := config.BackendAPI + "/v2"

	return api.NewClientWithResponses(backendServer, extraOptions...)
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

func newSDKHeaderEditorFn(_ context.Context, req *http.Request) error {
	sdk := struct {
		Name            string `json:"name"`
		SdkVersion      string `json:"sdkVersion"`
		LanguageVersion string `json:"languageVersion"`
	}{
		Name:            "Go SDK",
		SdkVersion:      Version,
		LanguageVersion: runtime.Version(),
	}

	marshaled, err := json.Marshal(sdk)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("X-Corbado-SDK", string(marshaled))

	return nil
}
