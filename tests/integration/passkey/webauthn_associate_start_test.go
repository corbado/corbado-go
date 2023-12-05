//go:build integration

package passkey_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebAuthnAssociateStart_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).Passkeys().AssociateStart(context.TODO(), api.WebAuthnAssociateStartReq{
		Origin:     "https://some.site.com",
		ClientInfo: util.ClientInfo("foobar", "127.0.0.1"),
	})
	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "associationToken: cannot be blank", servererror.GetValidationMessage(serverErr.Validation))
}
