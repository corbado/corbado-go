//go:build integration

package validation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestValidatePhoneNumber_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Validations().ValidatePhoneNumber(context.TODO(), api.ValidatePhoneNumberReq{
		PhoneNumber: "",
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "phoneNumber: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestValidatePhoneNumber_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Validations().ValidatePhoneNumber(context.TODO(), api.ValidatePhoneNumberReq{
		PhoneNumber: integration.CreateRandomTestPhoneNumber(t),
	})

	require.Nil(t, err)
	require.NotNil(t, rsp)
	assert.True(t, rsp.Data.IsValid)
	assert.Equal(t, api.PhoneNumberValidationResultValidationCodeValid, rsp.Data.ValidationCode)
	assert.NotNil(t, rsp.Data.PhoneNumber)
}
