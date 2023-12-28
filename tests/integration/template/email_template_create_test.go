//go:build integration

package template_test

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

func TestEmailTemplateCreate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Templates().CreateEmailTemplate(context.TODO(), api.EmailTemplateCreateReq{
		Lang:                     "invalid",
		Type:                     api.EmailTemplateCreateReqTypeEmailLink,
		Name:                     integration.CreateRandomTestName(t),
		Subject:                  "subject",
		PlainTextBody:            "plain body",
		HtmlTextTitle:            "html title",
		HtmlTextBody:             "html body",
		HtmlTextButton:           "html button",
		HtmlColorFont:            "#000000",
		HtmlColorBackgroundOuter: "#edf6ff",
		HtmlColorBackgroundInner: "#ffffff",
		HtmlColorButton:          "#1953ff",
		HtmlColorButtonFont:      "#ffffff",
		IsDefault:                false,
	})
	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "lang: must be in a valid format", servererror.GetValidationMessage(serverErr.Validation))
}

func TestEmailTemplateCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Templates().CreateEmailTemplate(context.TODO(), api.EmailTemplateCreateReq{
		Lang:                     "en",
		Type:                     api.EmailTemplateCreateReqTypeEmailLink,
		Name:                     integration.CreateRandomTestName(t),
		Subject:                  "subject",
		PlainTextBody:            "plain body",
		HtmlTextTitle:            "html title",
		HtmlTextBody:             "html body",
		HtmlTextButton:           "html button",
		HtmlColorFont:            "#000000",
		HtmlColorBackgroundOuter: "#edf6ff",
		HtmlColorBackgroundInner: "#ffffff",
		HtmlColorButton:          "#1953ff",
		HtmlColorButtonFont:      "#ffffff",
		IsDefault:                false,
	})
	require.Nil(t, err)
	require.NotNil(t, rsp)

	assert.NotEmpty(t, rsp.Data.EmailTemplateID)
}
