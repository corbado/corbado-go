//go:build integration

package passkey_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebAuthnCredentialUpdate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().CredentialUpdate(context.TODO(), "cre-12345678", api.WebAuthnCredentialReq{
		Status: api.WebAuthnCredentialReqStatusActive,
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "credentialID: does not exist", servererror.GetValidationMessage(serverErr.Validation))
}
