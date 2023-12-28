//go:build integration

package project_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/servererror"
	"github.com/corbado/corbado-go/pkg/util"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestProjectConfigUpdate_ValidationError(t *testing.T) {
	err := integration.SDK(t).Projects().ConfigUpdate(context.TODO(), api.ProjectConfigSaveReq{
		ExternalName: util.Ptr("a"),
	})

	require.NotNil(t, err)
	serverErr := corbado.AsServerError(err)
	require.NotNil(t, serverErr)

	assert.Equal(t, "externalName: the length must be between 2 and 255", servererror.GetValidationMessage(serverErr.Validation))
}

func TestProjectConfigUpdate_Success(t *testing.T) {
	newName := integration.CreateRandomTestName(t)
	err := integration.SDK(t).Projects().ConfigUpdate(context.TODO(), api.ProjectConfigSaveReq{
		ExternalName: util.Ptr(newName),
	})
	require.NoError(t, err)

	newCfg, err := integration.SDK(t).Projects().ConfigGet(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, newName, newCfg.Data.ExternalName)
}
