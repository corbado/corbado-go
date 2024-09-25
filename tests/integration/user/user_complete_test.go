//go:build integration

package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go/v2"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
	"github.com/corbado/corbado-go/v2/tests/integration"
)

func TestUserOperations(t *testing.T) {
	ctx := context.TODO()
	testUserID := ""

	t.Run("UserCreate", func(t *testing.T) {
		t.Run("ValidationError", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Create(ctx, api.UserCreateReq{})
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			require.Equal(t, "status: cannot be blank", serverErr.GetValidationMessage())
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Create(ctx, api.UserCreateReq{
				FullName: integration.CreateRandomTestName(t),
				Status:   api.UserStatusActive,
			})
			require.NotNil(t, rsp)
			require.NoError(t, err)

			testUserID = rsp.UserID
		})

		t.Run("SuccessActiveUser", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().CreateActiveByName(ctx, *integration.CreateRandomTestName(t))
			require.NotNil(t, rsp)
			require.NoError(t, err)
		})
	})

	t.Run("UserGet", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Get(ctx, "usr-123456789")
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			require.Equal(t, int32(400), serverErr.HTTPStatusCode)
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Get(ctx, testUserID)
			require.NotNil(t, rsp)
			require.NoError(t, err)
		})
	})

	t.Run("UserDelete", func(t *testing.T) {
		t.Run("ValidationError", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Delete(ctx, "usr-123456789")
			require.Nil(t, rsp)
			require.Error(t, err)

			serverErr := corbado.AsServerError(err)
			require.NotNil(t, serverErr)
			require.Equal(t, int32(400), serverErr.HTTPStatusCode)
			require.Equal(t, "userID: does not exist", serverErr.GetValidationMessage())
		})

		t.Run("Success", func(t *testing.T) {
			rsp, err := integration.SDK(t).Users().Delete(ctx, testUserID)
			require.NotNil(t, rsp)
			require.NoError(t, err)
		})
	})
}
