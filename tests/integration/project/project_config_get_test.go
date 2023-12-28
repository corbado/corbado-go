//go:build integration

package project_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestProjectConfigGet_AuthError(t *testing.T) {
	config, err := corbado.NewConfiguration("pro-12345678", "wrongsecret")
	require.NoError(t, err)
	config.BackendAPI = integration.GetBackendAPI(t)

	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	rsp, err := sdk.Projects().ConfigGet(context.TODO())

	require.Nil(t, rsp)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, int32(http.StatusUnauthorized), serverErr.HTTPStatusCode)
	assert.Equal(t, "login_error", serverErr.Type)
}

func TestProjectConfigGet_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Projects().ConfigGet(context.TODO())

	require.NoError(t, err)
	assert.Equal(t, integration.GetProjectID(t), rsp.Data.ProjectID)
}
