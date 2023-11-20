//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailLinkSend(t *testing.T) {
	sdk := initClient(t)
	rsp, err := sdk.EmailLinks().Send(context.TODO(), api.EmailLinkSendReq{
		Email:             createRandomTestEmail(t),
		Redirect:          "https://some.site.com/authenticate",
		TemplateName:      util.Ptr("default"),
		Create:            true,
		AdditionalPayload: util.Ptr("{}"),
		ClientInfo:        util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.EmailLinkID)
}
