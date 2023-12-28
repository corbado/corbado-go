//go:build integration

package validation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestValidateEmail_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Validations().ValidateEmail(context.TODO(), api.ValidateEmailReq{
		Email: "",
	})
	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "email: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestValidateEmail_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Validations().ValidateEmail(context.TODO(), api.ValidateEmailReq{
		Email: integration.CreateRandomTestEmail(t),
	})
	require.Nil(t, err)
	require.NotNil(t, rsp)

	assert.True(t, rsp.Data.IsValid)
	assert.Equal(t, api.EmailValidationResultValidationCodeValid, rsp.Data.ValidationCode)
	assert.NotNil(t, rsp.Data.Email)
}
