//go:build integration

package user

import (
	"context"
	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserCreate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Create(context.TODO(), api.UserCreateReq{})

	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "status: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestUserCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Create(context.TODO(), api.UserCreateReq{
		FullName: integration.CreateRandomTestName(t),
		Status:   api.UserStatusActive,
	})

	assert.NotNil(t, rsp)
	assert.NoError(t, err)
}

func TestActiveUserCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Users().CreateActiveByName(context.TODO(), *integration.CreateRandomTestName(t))

	assert.NotNil(t, rsp)
	assert.NoError(t, err)
}
