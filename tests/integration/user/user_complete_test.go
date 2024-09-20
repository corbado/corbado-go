//go:build integration

package user

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserOperations(t *testing.T) {
	ctx := context.TODO()
	testUserID := ""

	// Subtest for User Create functionality
	t.Run("UserCreate", func(t *testing.T) {
		t.Run("ValidationError", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Create(ctx, api.UserCreateReq{})
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			assert.Equal(t, "status: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Create(ctx, api.UserCreateReq{
				FullName: integration.CreateRandomTestName(t),
				Status:   api.UserStatusActive,
			})
			assert.NotNil(t, rsp)
			assert.NoError(t, err)

			testUserID = rsp.UserID
		})

		t.Run("CreateActiveUser_Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().CreateActiveByName(ctx, *integration.CreateRandomTestName(t))
			assert.NotNil(t, rsp)
			assert.NoError(t, err)
		})
	})

	// Subtest for User Get functionality
	t.Run("UserGet", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Get(ctx, "usr-123456789")
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			assert.Equal(t, int32(400), serverErr.HTTPStatusCode)
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Get(ctx, testUserID)
			require.NotNil(t, rsp)
			require.NoError(t, err)
		})
	})

	// Subtest for User Delete functionality
	t.Run("UserDelete", func(t *testing.T) {
		t.Run("ValidationError", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Delete(ctx, "usr-123456789")
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			assert.Equal(t, int32(400), serverErr.HTTPStatusCode)
			assert.Equal(t, "userID: does not exist", servererror.GetValidationMessage(serverErr.Validation))
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Delete(ctx, testUserID)
			require.NotNil(t, rsp)
			require.NoError(t, err)
		})
	})
}
