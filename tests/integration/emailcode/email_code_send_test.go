//go:build integration

package emailcode_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailCodeSend(t *testing.T) {
	rsp, err := integration.SDK(t).EmailCodes().Send(context.TODO(), api.EmailCodeSendReq{
		Email:             integration.CreateRandomTestEmail(t),
		TemplateName:      util.Ptr("default"),
		Create:            true,
		AdditionalPayload: util.Ptr("{}"),
		ClientInfo:        util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.EmailCodeID)
}
