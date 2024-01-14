//go:build integration

package smsotp_test

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

func TestSmsOTPValidate_ValidationErrorEmptyCode(t *testing.T) {
	_, err := integration.SDK(t).SmsOTPs().Validate(context.TODO(), "sms-123456789", api.SmsCodeValidateReq{
		SmsCode: "",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "smsCode: cannot be blank")
}

func TestSmsOTPValidate_ValidationErrorInvalidCode(t *testing.T) {
	_, err := integration.SDK(t).SmsOTPs().Validate(context.TODO(), "sms-123456789", api.SmsCodeValidateReq{
		SmsCode: "1",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "smsCode: the length must be exactly 6")
}

func TestSmsOTPValidate_ValidationErrorInvalidID(t *testing.T) {
	_, err := integration.SDK(t).SmsOTPs().Validate(context.TODO(), "sms-123456789", api.SmsCodeValidateReq{
		SmsCode: "123456",
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, int32(404), serverErr.HTTPStatusCode)
}

func TestSmsOTPValidate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).SmsOTPs().Send(context.TODO(), api.SmsCodeSendReq{
		PhoneNumber:  integration.CreateRandomTestPhoneNumber(t),
		TemplateName: util.Ptr("default"),
		Create:       true,
		ClientInfo:   util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.SmsCodeID)

	_, err = integration.SDK(t).SmsOTPs().Validate(context.TODO(), rsp.Data.SmsCodeID, api.SmsCodeValidateReq{
		SmsCode: "150919",
	})

	require.NoError(t, err)
}
