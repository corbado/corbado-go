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

func TestWebAuthnAuthenticateStart_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().AuthenticateStart(context.TODO(), api.WebAuthnAuthenticateStartReq{
		Username:   "",
		ClientInfo: *util.ClientInfo("foobar", "127.0.0.1"),
	})
	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "username: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}
