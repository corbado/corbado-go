//go:build integration

package authtoken_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/servererror"
	"github.com/corbado/corbado-go/pkg/util"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestAuthTokenValidate_ValidationError(t *testing.T) {
	rsp, err := integration.SDK(t).AuthTokens().Validate(context.TODO(), api.AuthTokenValidateReq{
		Token:      "invalid",
		ClientInfo: *util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "token: the length must be exactly 64", servererror.GetValidationMessage(serverErr.Validation))
}

func TestAuthTokenValidate_NotExists(t *testing.T) {
	rsp, err := integration.SDK(t).AuthTokens().Validate(context.TODO(), api.AuthTokenValidateReq{
		Token:      strings.Repeat("a", 64),
		ClientInfo: *util.ClientInfo("foobar", "127.0.0.1"),
	})

	require.Nil(t, rsp)
	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, int32(http.StatusNotFound), serverErr.HTTPStatusCode)
	require.NotNil(t, serverErr.Details)
	assert.Equal(t, "Session doesn't exist", *serverErr.Details)
}
