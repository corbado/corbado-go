//go:build integration

package emailcode_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/util"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestEmailCodeSend(t *testing.T) {
	rsp, err := integration.SDK(t).EmailOTPs().Send(context.TODO(), api.EmailCodeSendReq{
		Email:             integration.CreateRandomTestEmail(t),
		TemplateName:      util.Ptr("default"),
		Create:            true,
		AdditionalPayload: util.Ptr("{}"),
		ClientInfo:        util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Data.EmailCodeID)
}
