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

func TestSmsOTPSend_ValidationError(t *testing.T) {
	_, err := integration.SDK(t).SmsOTPs().Send(context.TODO(), api.SmsCodeSendReq{
		PhoneNumber:  "",
		TemplateName: util.Ptr("default"),
		Create:       true,
		ClientInfo:   util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.Error(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "phoneNumber: cannot be blank")
}

func TestSmsOTPSend_Success(t *testing.T) {
	rsp, err := integration.SDK(t).SmsOTPs().Send(context.TODO(), api.SmsCodeSendReq{
		PhoneNumber:  integration.CreateRandomTestPhoneNumber(t),
		TemplateName: util.Ptr("default"),
		Create:       true,
		ClientInfo:   util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.SmsCodeID)
}
