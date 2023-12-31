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

func TestUserCreate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Create(context.TODO(), api.UserCreateReq{
		Name:  "",
		Email: util.Ptr(""),
	})

	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "name: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestUserCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Create(context.TODO(), api.UserCreateReq{
		Name:  integration.CreateRandomTestName(t),
		Email: util.Ptr(integration.CreateRandomTestEmail(t)),
	})

	assert.NotNil(t, rsp)
	assert.NoError(t, err)
}
