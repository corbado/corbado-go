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
	"github.com/corbado/corbado-go/pkg/util"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestWebAuthnCredentialList_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().CredentialList(context.TODO(), &api.WebAuthnCredentialListParams{
		Sort: util.Ptr("foo:bar"),
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "sort: Invalid order direction 'bar'", servererror.GetValidationMessage(serverErr.Validation))
}

func TestWebAuthnCredentialList_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().CredentialList(context.TODO(), nil)
	require.NoError(t, err)

	assert.True(t, len(rsp.Rows) == 0)
}
