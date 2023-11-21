//go:build integration

package user_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserList(t *testing.T) {
	// send email link first so that we have at least one user
	_, err := integration.SDK(t).EmailLinks().Send(context.TODO(), api.EmailLinkSendReq{
		Email:             integration.CreateRandomTestEmail(t),
		Redirect:          "https://some.site.com/authenticate",
		TemplateName:      util.Ptr("default"),
		Create:            true,
		AdditionalPayload: util.Ptr("{}"),
		ClientInfo:        util.ClientInfo("foobar", "127.0.0.1"),
	})
	require.NoError(t, err)

	usersRsp, err := integration.SDK(t).Users().List(context.TODO(), nil)
	require.NoError(t, err)

	assert.True(t, len(usersRsp.Data.Users) > 0)
}
