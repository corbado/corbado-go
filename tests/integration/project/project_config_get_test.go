//go:build integration

package project_test

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectConfigGet(t *testing.T) {
	rsp, err := integration.SDK(t).Projects().ConfigGet(context.TODO())

	require.NoError(t, err)
	assert.Equal(t, integration.GetProjectID(t), rsp.Data.ProjectID)
}
