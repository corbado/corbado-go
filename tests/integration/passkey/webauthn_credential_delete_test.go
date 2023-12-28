//go:build integration

package passkey_test

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

func TestWebAuthnCredentialDelete_ValidationError(t *testing.T) {
	err := integration.SDK(t).Passkeys().CredentialDelete(context.TODO(), "usr-12345678", "cre-12345678", api.EmptyReq{})

	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "credentialID: does not exist", servererror.GetValidationMessage(serverErr.Validation))
}
