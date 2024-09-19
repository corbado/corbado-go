//go:build integration

package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestUserGet_NotFound(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Get(context.TODO(), "usr-123456789")
	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, int32(400), serverErr.HTTPStatusCode)
}

func TestUserGet_Success(t *testing.T) {
	userID := integration.CreateUser(t)

	rsp, err := integration.SDK(t).Users().Get(context.TODO(), userID)
	require.NotNil(t, rsp)
	require.NoError(t, err)
}
