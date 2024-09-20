package session

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/MicahParks/keyfunc"
	"github.com/corbado/corbado-go/pkg/entities"
	"github.com/corbado/corbado-go/pkg/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/corbado/corbado-go/internal/services/session"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// Helper function to generate JWTs
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

// provideJWTs returns a slice of test cases for JWT validation
func provideJWTs(privateKey *rsa.PrivateKey) [][]interface{} {
	return [][]interface{}{
		{
			// JWT with invalid format
			"invalid",
			false,
		},
		{
			// JWT with invalid signature
			"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImtpZDEyMyJ9.eyJpc3MiOiJodHRwczovL2F1dGguYWNtZS5jb20iLCJpYXQiOjE3MjY0OTE4MDcsImV4cCI6MTcyNjQ5MTkwNywibmJmIjoxNzI2NDkxNzA3LCJzdWIiOiJ1c3ItMTIzNDU2Nzg5MCIsIm5hbWUiOiJuYW1lIiwiZW1haWwiOiJlbWFpbCIsInBob25lX251bWJlciI6InBob25lTnVtYmVyIiwib3JpZyI6Im9yaWcifQ.invalid",
			false,
		},
		{
			// Not before (nbf) in future
			generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Add(100*time.Second).Unix(), privateKey),
			false,
		},
		{
			// Expired (exp)
			generateJWT("https://auth.acme.com", time.Now().Add(-100*time.Second).Unix(), time.Now().Add(-100*time.Second).Unix(), privateKey),
			false,
		},
		{
			// Invalid issuer (iss)
			generateJWT("https://invalid.com", time.Now().Add(100*time.Second).Unix(), time.Now().Add(100*time.Second).Unix(), privateKey),
			false,
		},
		{
			// Success
			generateJWT("https://auth.acme.com", time.Now().Add(100*time.Second).Unix(), time.Now().Unix(), privateKey),
			true, // Success case
		},
	}
}

// CreateSession mocks the JWKS endpoint and creates a new SessionService
func createSession(config *session.Config) (*session.Impl, error) {
	// Read JWKS data from a file
	// Get the absolute path relative to the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	// Construct the absolute path to the private key file
	filePath := filepath.Join(workingDir, "../testdata/jwks.json")

	jwksData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read JWKS file")
	}

	// Create an HTTP mock server for the JWKS endpoint
	mockServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(jwksData)
	})

	// Start the mock server
	server := &http.Server{Addr: "localhost:8081", Handler: mockServer}
	go func() {
		_ = server.ListenAndServe()
	}()

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
		ResponseExtractor: func(ctx context.Context, resp *http.Response) (json.RawMessage, error) {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
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

	// JWKS URI should be the mock server URL
	jwksURI := "http://localhost:8081"

	jwks, err := keyfunc.Get(jwksURI, options)
	if err != nil {
		return nil, errors.New("failed to create JWKS")
	}

	// Create the SessionService implementation with the mocked JWKS
	sessionService := &session.Impl{
		Client: &api.ClientWithResponses{},
		Config: config,
		Jwks:   jwks,
	}

	return sessionService, nil
}

func TestValidateToken(t *testing.T) {
	// Create a mock session using the createSession function
	config := &session.Config{
		ProjectID:            "test-project-id",
		JwksURI:              "http://localhost:8081",
		JWTIssuer:            "https://auth.acme.com",
		JWKSRefreshInterval:  0,
		JWKSRefreshRateLimit: 0,
		JWKSRefreshTimeout:   0,
	}
	session, err := createSession(config)
	assert.NoError(t, err)

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	// Construct the absolute path to the private key file
	filePath := filepath.Join(workingDir, "../testdata/privateKey.pem")

	privateKeyFile, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	jwtCases := provideJWTs(privateKey)

	for _, testCase := range jwtCases {
		shortSession := testCase[0].(string)
		success := testCase[1].(bool)

		t.Run(shortSession, func(t *testing.T) {
			var exception error
			var user *entities.User

			// Try to validate the token
			user, err := session.ValidateToken(shortSession)
			if err != nil {
				exception = err
			}

			if success {
				assert.NotNil(t, user, "Expected a user but got nil")
				assert.Equal(t, "usr-1234567890", user.UserID, "User ID should be 'usr-1234567890'")
			} else {
				assert.NotNil(t, exception, "Expected an exception but got nil")
				var validationErr *jwt.ValidationError
				assert.ErrorAs(t, exception, &validationErr, "Expected a ValidationError")
				log.Println(validationErr.Errors)

				switch {
				case validationErr.Errors&jwt.ValidationErrorMalformed != 0:
					assert.Equal(t, jwt.ValidationErrorMalformed, validationErr.Errors&jwt.ValidationErrorMalformed, "Expected a malformed token error")
				case validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
					assert.Equal(t, jwt.ValidationErrorSignatureInvalid, validationErr.Errors&jwt.ValidationErrorSignatureInvalid, "Expected an invalid signature error")
				case validationErr.Errors&jwt.ValidationErrorNotValidYet != 0:
					assert.Equal(t, jwt.ValidationErrorNotValidYet, validationErr.Errors&jwt.ValidationErrorNotValidYet, "Expected a 'not valid yet' error")
				case validationErr.Errors&jwt.ValidationErrorExpired != 0:
					assert.Equal(t, jwt.ValidationErrorExpired, validationErr.Errors&jwt.ValidationErrorExpired, "Expected an expired token error")
				case validationErr.Errors&jwt.ValidationErrorIssuer != 0:
					assert.Equal(t, jwt.ValidationErrorIssuer, validationErr.Errors&jwt.ValidationErrorIssuer, "Expected an invalid issuer error")
				default:
					t.Fatalf("Unexpected validation error: 0x%x", validationErr.Errors)
				}

			}
		})
	}
}
