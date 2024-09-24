package session

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go/pkg/logger"
	"github.com/corbado/corbado-go/pkg/validationerror"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/corbado/corbado-go/internal/services/session"
	"github.com/corbado/corbado-go/pkg/generated/api"
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
func newSession() (*session.Impl, error) {
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

	server := &http.Server{Addr: "localhost:8081", Handler: mockServer} // nolint:gosec
	go func() {
		_ = server.ListenAndServe()
	}()

	// Config
	config := &session.Config{
		ProjectID:            "test-project-id",
		JwksURI:              "http://localhost:8081",
		JWTIssuer:            "https://auth.acme.com",
		JWKSRefreshInterval:  0,
		JWKSRefreshRateLimit: 0,
		JWKSRefreshTimeout:   0,
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

	jwks, err := keyfunc.Get("http://localhost:8081", options)
	if err != nil {
		return nil, err
	}

	return &session.Impl{
		Client: &api.ClientWithResponses{},
		Config: config,
		Jwks:   jwks,
	}, nil
}

func TestValidateToken(t *testing.T) {
	validPrivateKey, err := generatePrivateKey("validPrivateKey.pem")
	require.NoError(t, err)

	invalidPrivateKey, err := generatePrivateKey("invalidPrivateKey.pem")
	require.NoError(t, err)

	tests := []struct {
		name                string
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
			shortSession:        "invalid",
			validationErrorCode: validationerror.CodeJWTInvalidData,
			success:             false,
		},
		{
			name: "JWT with invalid signature",
			// nolint:lll
			shortSession:        "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImtpZDEyMyJ9.eyJpc3MiOiJodHRwczovL2F1dGguYWNtZS5jb20iLCJpYXQiOjE3MjY0OTE4MDcsImV4cCI6MTcyNjQ5MTkwNywibmJmIjoxNzI2NDkxNzA3LCJzdWIiOiJ1c3ItMTIzNDU2Nzg5MCIsIm5hbWUiOiJuYW1lIiwiZW1haWwiOiJlbWFpbCIsInBob25lX251bWJlciI6InBob25lTnVtYmVyIiwib3JpZyI6Im9yaWcifQ.invalid",
			validationErrorCode: validationerror.CodeJWTInvalidSignature,
			success:             false,
		},
		{
			name:                "JWT with invalid private key signed",
			shortSession:        generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), invalidPrivateKey),
			validationErrorCode: validationerror.CodeJWTInvalidSignature,
			success:             false,
		},
		{
			name:                "Not before (nbf) in future",
			shortSession:        generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Add(100*time.Second).Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTBefore,
			success:             false,
		},
		{
			name:                "Expired (exp)",
			shortSession:        generateJWT("https://auth.acme.com", time.Now().Add(-100*time.Second).Unix(), time.Now().Add(-100*time.Second).Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTExpired,
			success:             false,
		},
		{
			name:                "Invalid issuer (iss)",
			shortSession:        generateJWT("https://invalid.com", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			validationErrorCode: validationerror.CodeJWTIssuerMismatch,
			success:             false,
		},
		{
			name:         "Success",
			shortSession: generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), validPrivateKey),
			success:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionSvc, err := newSession()
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
