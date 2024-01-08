//go:build integration

package emailotp_test

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

func TestEmailOTPValidate_ValidationErrorEmptyCode(t *testing.T) {
	_, err := integration.SDK(t).EmailOTPs().Validate(context.TODO(), "emc-123456789", api.EmailCodeValidateReq{
		Code: "",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "code: cannot be blank")
}

func TestEmailOTPValidate_ValidationErrorInvalidCode(t *testing.T) {
	_, err := integration.SDK(t).EmailOTPs().Validate(context.TODO(), "emc-123456789", api.EmailCodeValidateReq{
		Code: "1",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "code: the length must be exactly 6")
}

func TestEmailOTPValidate_ValidationErrorInvalidID(t *testing.T) {
	_, err := integration.SDK(t).EmailOTPs().Validate(context.TODO(), "emc-123456789", api.EmailCodeValidateReq{
		Code: "123456",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, int32(404), serverErr.HTTPStatusCode)
}

func TestEmailOTPValidate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).EmailOTPs().Send(context.TODO(), api.EmailCodeSendReq{
		Email:             integration.CreateRandomTestEmail(t),
		TemplateName:      util.Ptr("default"),
		Create:            true,
		AdditionalPayload: util.Ptr("{}"),
		ClientInfo:        util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.EmailCodeID)

	_, err = integration.SDK(t).EmailOTPs().Validate(context.TODO(), rsp.Data.EmailCodeID, api.EmailCodeValidateReq{
		Code: "150919",
	})

	require.NoError(t, err)
}
