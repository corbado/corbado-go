//go:build integration

package identifier

import (
	"context"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_DeleteSuccessful(t *testing.T) {
	// Create Identifier
	identifierId, userId, _ := integration.CreateIdentifier(t)

	// List Identifiers
	initialIdentifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	if err != nil {
		panic(err)
	}

	// Delete identifiers
	identifier, err := integration.SDK(t).Identifiers().Delete(context.Background(), userId, identifierId)

	// List Identifiers after deletion
	identifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, identifier)
	assert.Equal(t, len(identifierList.Identifiers), len(initialIdentifierList.Identifiers)-1)
}

func TestIdentifier_DeleteFailure(t *testing.T) {
	// Create User
	userID := integration.CreateUser(t)
	// Delete unknown identifier
	identifier, err := integration.SDK(t).Identifiers().Delete(context.Background(), userID, "ide-123456")

	assert.Error(t, err)
	assert.Nil(t, identifier)
}
