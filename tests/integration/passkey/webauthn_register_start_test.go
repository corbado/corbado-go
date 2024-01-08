//go:build integration

package passkey_test

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

func TestWebAuthnRegisterStart_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().RegisterStart(context.TODO(), api.WebAuthnRegisterStartReq{
		Username:   "",
		ClientInfo: *util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Equal(t, "username: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}

func TestWebAuthnRegisterStart_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().RegisterStart(context.TODO(), api.WebAuthnRegisterStartReq{
		Username:     integration.CreateRandomTestEmail(t),
		UserFullName: util.Ptr(integration.CreateRandomTestName(t)),
		ClientInfo:   *util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.PublicKeyCredentialCreationOptions)
	assert.Equal(t, "success", string(rsp.Status))
}
