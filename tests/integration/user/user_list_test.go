//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/pkg/util"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestUserList_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Users().List(context.TODO(), &api.UserListParams{
		Sort: util.Ptr("foo:bar"),
	})

	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "sort: Invalid order direction 'bar'", servererror.GetValidationMessage(serverErr.Validation))
}

func TestUserList_Success(t *testing.T) {
	userID := integration.CreateUser(t)

	rsp, err := integration.SDK(t).Users().List(context.TODO(), &api.UserListParams{
		Sort: util.Ptr("created:desc"),
	})

	require.NoError(t, err)

	found := false
	for _, user := range rsp.Data.Users {
		if user.ID == userID {
			found = true
			break
		}
	}

	assert.True(t, found)
}
