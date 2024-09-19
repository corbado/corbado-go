//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestUserDelete_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Delete(context.TODO(), "usr-123456789")
	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, int32(400), serverErr.HTTPStatusCode)
	assert.Equal(t, "userID: does not exist", servererror.GetValidationMessage(serverErr.Validation))
}

func TestUserDelete_Success(t *testing.T) {
	userID := integration.CreateUser(t)

	rsp, err := integration.SDK(t).Users().Delete(context.TODO(), userID)
	require.NotNil(t, rsp)
	require.NoError(t, err)
}
