//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/generated/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestUserExists_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Exists(context.TODO(), api.UserExistsReq{})

	require.Nil(t, rsp)
	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "loginIdentifier: cannot be blank")
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "loginIdentifierType: cannot be blank")
}

func TestUserExists_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Users().Exists(context.TODO(), api.UserExistsReq{
		LoginIdentifierType: common.LoginIdentifierTypeEmail,
		LoginIdentifier:     "idontexist@somemail.com",
	})

	assert.NotNil(t, rsp)
	assert.NoError(t, err)
	assert.False(t, rsp.Exists)
}
