package session

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go/v2/pkg/logger"
	"github.com/corbado/corbado-go/v2/pkg/validationerror"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/corbado/corbado-go/v2/internal/services/session"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
)

func generateJWT(iss string, exp, nbf int64, privateKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":          iss,
		"iat":          time.Now().Unix(),
		"exp":          exp,
		"nbf":          nbf,
		"sub":          "usr-1234567890",
		"name":         "name",
		"email":        "email",
		"phone_number": "phoneNumber",
		"orig":         "orig",
	})

	token.Header["kid"] = "kid123"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func generatePrivateKey(filename string) (*rsa.PrivateKey, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	privateKeyFile, err := os.ReadFile(filepath.Join(workingDir, "../testdata/"+filename))
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
}

// newSession mocks the JWKS endpoint and creates a new session service
func newSession(issuer string) (*session.Impl, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	jwksData, err := os.ReadFile(filepath.Join(workingDir, "../testdata/jwks.json"))
	if err != nil {
		return nil, err
	}

	// Create an HTTP mock server for the JWKS endpoint
	mockServer := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(jwksData); err != nil {
			panic(err)
		}
	})

	server := httptest.NewServer(mockServer)
	defer server.Close()

	// Config
	config := &session.Config{
		ProjectID: "pro-1",
		JwksURI:   server.URL,
		JWTIssuer: issuer,
	}

	// Create a new JWKS instance using the mock JWKS server
	options := keyfunc.Options{
		RequestFactory: func(ctx context.Context, urlAddress string) (*http.Request, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddress, nil)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Corbado-ProjectID", config.ProjectID)

			return req, nil
		},
		ResponseExtractor: func(_ context.Context, resp *http.Response) (json.RawMessage, error) {
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return body, nil
		},
		RefreshErrorHandler: func(err error) {
			logger.Error("Error refreshing JWKS: %s", err.Error())
		},
		RefreshInterval:   config.JWKSRefreshInterval,
		RefreshRateLimit:  config.JWKSRefreshRateLimit,
		RefreshTimeout:    config.JWKSRefreshTimeout,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(server.URL, options)
	if err != nil {
		return nil, err
	}

	return &session.Impl{
		Client: &api.ClientWithResponses{},
		Config: config,
		Jwks:   jwks,
	}, nil
}

// nolint:funlen
func TestValidateToken(t *testing.T) {
	validPrivateKey, err := generatePrivateKey("validPrivateKey.pem")
	require.NoError(t, err)

	invalidPrivateKey, err := generatePrivateKey("invalidPrivateKey.pem")
	require.NoError(t, err)

	tests := []struct {
		name                string
		issuer              string
		shortSession        string
		validationErrorCode validationerror.Code
		success             bool
	}{
		{
			name:         "Empty JWT",
			shortSession: "",
			success:      false,
		},
		{
			name:                "JWT with invalid format",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        "invalid",
			validationErrorCode: validationerror.CodeJWTInvalidData,
			success:             false,
		},
		{
			name:                "JWT with invalid signature",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImtpZDEyMyJ9.eyJpc3MiOiJodHRwczovL2F1dGguYWNtZS5jb20iLCJpYXQiOjE3MjY0OTE4MDcsImV4cCI6MTcyNjQ5MTkwNywibmJmIjoxNzI2NDkxNzA3LCJzdWIiOiJ1c3ItMTIzNDU2Nzg5MCIsIm5hbWUiOiJuYW1lIiwiZW1haWwiOiJlbWFpbCIsInBob25lX251bWJlciI6InBob25lTnVtYmVyIiwib3JpZyI6Im9yaWcifQ.invalid", // nolint:lll
			validationErrorCode: validationerror.CodeJWTInvalidSignature,
			success:             false,
		},
		{
			name:                "JWT with invalid private key signed",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        generateJWT("https://pro-1.frontendapi.cloud.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), invalidPrivateKey),
			validationErrorCode: validationerror.CodeJWTInvalidSignature,
			success:             false,
		},
		{
			name:                "Not before (nbf) in future",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        generateJWT("https://pro-1.frontendapi.cloud.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Add(100*time.Second).Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTBefore,
			success:             false,
		},
		{
			name:                "Expired (exp)",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        generateJWT("https://pro-1.frontendapi.cloud.corbado.io", time.Now().Add(-100*time.Second).Unix(), time.Now().Add(-100*time.Second).Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTExpired,
			success:             false,
		},
		{
			name:                "Empty issuer (iss)",
			issuer:              "https://pro-1.frontendapi.corbado.io",
			shortSession:        generateJWT("", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTIssuerEmpty,
			success:             false,
		},
		{
			name:                "Invalid issuer 1 (iss)",
			issuer:              "https://pro-1.frontendapi.corbado.io",
			shortSession:        generateJWT("https://pro-2.frontendapi.cloud.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTIssuerMismatch,
			success:             false,
		},
		{
			name:                "Invalid issuer 2 (iss)",
			issuer:              "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession:        generateJWT("https://pro-2.frontendapi.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTIssuerMismatch,
			success:             false,
		},
		{
			name:         "Success with old Frontend API URL in JWT",
			issuer:       "https://pro-1.frontendapi.cloud.corbado.io",
			shortSession: generateJWT("https://pro-1.frontendapi.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			success:      true,
		},
		{
			name:         "Success with old Frontend API URL in config",
			issuer:       "https://pro-1.frontendapi.corbado.io",
			shortSession: generateJWT("https://pro-1.frontendapi.cloud.corbado.io", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			success:      true,
		},
		{
			name:         "Success with CNAME",
			issuer:       "https://auth.acme.com",
			shortSession: generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			success:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionSvc, err := newSession(test.issuer)
			require.NoError(t, err)

			user, err := sessionSvc.ValidateToken(test.shortSession)

			if test.success {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "usr-1234567890", user.UserID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)

				if test.validationErrorCode > 0 {
					var validationErr *validationerror.ValidationError
					assert.ErrorAs(t, err, &validationErr)
					assert.Equal(t, test.validationErrorCode, validationErr.Code)
				}
			}
		})
	}
}
