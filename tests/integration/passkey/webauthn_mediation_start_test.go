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

func TestWebAuthnMediationStart_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().MediationStart(context.TODO(), api.WebAuthnMediationStartReq{})
	require.Nil(t, rsp)
	require.NotNil(t, err)

	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)
	assert.Contains(t, servererror.GetValidationMessage(serverErr.Validation), "userAgent: cannot be blank")
}

func TestWebAuthnMediationStart_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().MediationStart(context.TODO(), api.WebAuthnMediationStartReq{
		ClientInfo: *util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, rsp.Challenge)
}
