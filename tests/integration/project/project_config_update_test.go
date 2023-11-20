//go:build integration

package project_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/sdk/entity/api"
	"github.com/corbado/corbado-go/pkg/sdk/util"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectConfigUpdate(t *testing.T) {
	newName := integration.CreateRandomTestName(t)
	err := integration.SDK(t).Projects().ConfigUpdate(context.TODO(), api.ProjectConfigSaveReq{
		ExternalName: util.Ptr(newName),
	})
	require.NoError(t, err)

	newCfg, err := integration.SDK(t).Projects().ConfigGet(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, newName, newCfg.Data.ExternalName)
}
