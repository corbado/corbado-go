//go:build integration

package session_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/entities"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/golang-jwt/jwt/v4"
)

func TestSession_ValidateToken_AuthError(t *testing.T) {
	// Invalid secret to trigger authentication error
	config, err := corbado.NewConfig("pro-12345678", "corbado1_wrongsecret")
	require.NoError(t, err)

	config.BackendAPI = integration.GetBackendAPI(t)
	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	// Test session validation failure
	shortSession := "invalid.jwt.token"
	_, err = sdk.Sessions().ValidateToken(context.TODO(), shortSession)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, int32(http.StatusUnauthorized), serverErr.HTTPStatusCode)
	assert.Equal(t, "login_error", serverErr.Type)
}

func TestSession_ValidateToken_Success(t *testing.T) {
	// Valid session token and proper setup
	rsp, err := integration.SDK(t).Sessions().ValidateToken(context.TODO(), "valid.jwt.token")
	require.NoError(t, err)

	// Check the returned user information
	assert.Equal(t, "12345", rsp.UserID)
	assert.Equal(t, "Test Name", rsp.FullName)
}

func TestSession_ValidateToken_IssuerMismatch(t *testing.T) {
	// Setup SDK with valid config
	sdk := integration.SDK(t)

	// Mock JWT with invalid issuer
	mockClaims := &entities.Claims{
		Issuer:  "https://invalid-issuer.com",
		Subject: "12345",
		Name:    "Test Name",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mockClaims)
	tokenString, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	// Attempt to validate the token with the wrong issuer
	_, err = sdk.Sessions().ValidateToken(context.TODO(), tokenString)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "issuer_mismatch", serverErr.Type)
}

func TestSession_ValidateToken_JWTExpired(t *testing.T) {
	// Setup SDK
	sdk := integration.SDK(t)

	// Mock an expired JWT
	mockClaims := &entities.Claims{
		Issuer:  "https://auth.example.com",
		Subject: "12345",
		Name:    "Test Name",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mockClaims)
	tokenString, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	// Attempt to validate the expired token
	_, err = sdk.Sessions().ValidateToken(context.TODO(), tokenString)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "token_expired", serverErr.Type)
}
