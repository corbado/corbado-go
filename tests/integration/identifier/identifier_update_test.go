//go:build integration

package identifier

import (
	"context"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_UpdateSuccessful(t *testing.T) {
	// Create Identifier
	identifierId, userId, _ := integration.CreateIdentifier(t)

	// List Identifiers
	initialIdentifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	assert.Nil(t, err)
	assert.Equal(t, initialIdentifierList.Identifiers[0].Status, api.IdentifierStatusVerified)

	if err != nil {
		panic(err)
	}

	// update identifier
	identifier, err := integration.SDK(t).Identifiers().UpdateIdentifier(context.Background(), userId, identifierId, api.IdentifierUpdateReq{Status: api.IdentifierStatusPending})

	if err != nil {
		panic(err)
	}

	// List Identifiers after update
	identifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, identifier)
	assert.Equal(t, identifierList.Identifiers[0].Status, api.IdentifierStatusPending)
}

func TestIdentifier_UpdateStatusSuccess(t *testing.T) {
	// Create Identifier
	identifierId, userId, _ := integration.CreateIdentifier(t)

	// List Identifiers
	initialIdentifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	assert.Nil(t, err)
	assert.Equal(t, initialIdentifierList.Identifiers[0].Status, api.IdentifierStatusVerified)

	if err != nil {
		panic(err)
	}

	// update identifier
	identifier, err := integration.SDK(t).Identifiers().UpdateStatus(context.Background(), userId, identifierId, api.IdentifierStatusPending)

	if err != nil {
		panic(err)
	}

	// List Identifiers after update
	identifierList, err := integration.SDK(t).Identifiers().ListByUserID(context.Background(), userId, "", 1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, identifier)
	assert.Equal(t, identifierList.Identifiers[0].Status, api.IdentifierStatusPending)
}
