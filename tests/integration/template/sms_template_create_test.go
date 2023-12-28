//go:build integration

package template_test

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

func TestSMSTemplateCreate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Templates().CreateSMSTemplate(context.TODO(), api.SmsTemplateCreateReq{
		IsDefault: false,
		Name:      integration.CreateRandomTestName(t),
		TextPlain: "",
		Type:      api.SmsTemplateCreateReqTypeSmsCode,
	})
	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "textPlain: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestSMSTemplateCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Templates().CreateSMSTemplate(context.TODO(), api.SmsTemplateCreateReq{
		IsDefault: false,
		Name:      integration.CreateRandomTestName(t),
		TextPlain: "text plain",
		Type:      api.SmsTemplateCreateReqTypeSmsCode,
	})
	require.Nil(t, err)
	require.NotNil(t, rsp)

	assert.NotEmpty(t, rsp.Data.SmsTemplateID)
}
