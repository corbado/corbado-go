//go:build integration

package session_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionConfigGet_AuthError(t *testing.T) {
	config, err := corbado.NewConfig("pro-12345678", "wrongsecret")
	require.NoError(t, err)
	config.BackendAPI = integration.GetBackendAPI(t)

	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	rsp, err := sdk.Sessions().ConfigGet(context.TODO(), nil)

	require.Nil(t, rsp)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, int32(http.StatusUnauthorized), serverErr.HTTPStatusCode)
	assert.Equal(t, "login_error", serverErr.Type)
}

func TestSessionConfigGet_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Sessions().ConfigGet(context.TODO(), nil)

	require.NoError(t, err)
	assert.Equal(t, integration.GetProjectID(t), rsp.Data.ProjectID)
}
