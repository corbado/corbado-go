//go:build integration

package identifier

import (
	"context"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_List(t *testing.T) {
	// List identifiers
	initialList, err := integration.SDK(t).Identifiers().List(context.Background(), []string{}, "", 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, initialList)

	// Create new identifier
	integration.CreateIdentifier(t)

	// List identifiers containing new identifier
	newList, err := integration.SDK(t).Identifiers().List(context.Background(), []string{}, "", 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, initialList)
	assert.Equal(t, len(newList.Identifiers), len(initialList.Identifiers)+1)
}

func TestIdentifier_ListByValueAndType(t *testing.T) {
	// Create new identifier
	_, _, email := integration.CreateIdentifier(t)

	// List identifiers
	initialList, err := integration.SDK(t).Identifiers().ListByValueAndType(context.Background(), email, "email", "", 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, initialList)
	assert.Len(t, initialList.Identifiers, 1)
	assert.Equal(t, initialList.Identifiers[0].Value, email)
}

func TestIdentifier_ListByUserIDAndType(t *testing.T) {
	// Create new identifier
	_, userId, email := integration.CreateIdentifier(t)

	// List identifiers
	initialList, err := integration.SDK(t).Identifiers().ListByUserIDAndType(context.Background(), userId, "email", "", 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, initialList)
	assert.Len(t, initialList.Identifiers, 1)
	assert.Equal(t, initialList.Identifiers[0].Value, email)
}

func TestIdentifier_ListByUserID(t *testing.T) {
	// Create new identifier
	_, userId, email := integration.CreateIdentifier(t)

	// List identifiers
	initialList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, initialList)
	assert.Len(t, initialList.Identifiers, 1)
	assert.Equal(t, initialList.Identifiers[0].Value, email)
}
