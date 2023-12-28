//go:build integration

package project_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestProjectSecretCreate_Success(t *testing.T) {
	rsp, err := integration.SDK(t).Projects().CreateSecret(context.TODO(), api.ProjectSecretCreateReq{})
	require.Nil(t, err)
	require.NotNil(t, rsp)

	require.NotNil(t, rsp.Secret)
	assert.NotEmpty(t, *rsp.Secret)
	assert.NotEmpty(t, rsp.Hint)
}
