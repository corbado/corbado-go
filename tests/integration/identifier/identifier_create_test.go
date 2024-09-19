//go:build integration

package identifier

import (
	"context"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_CreateSuccesful(t *testing.T) {
	// Create User
	userID := integration.CreateUser(t)

	email := integration.CreateRandomTestEmail(t)
	// List identifiers
	identifier, err := integration.SDK(t).Identifiers().Create(context.Background(), userID, api.IdentifierCreateReq{IdentifierValue: email, IdentifierType: "email", Status: "verified"})
	assert.NoError(t, err)
	assert.NotNil(t, identifier)
	assert.Equal(t, identifier.Value, email)
	assert.Equal(t, identifier.Status, api.IdentifierStatus("verified"))
	assert.Equal(t, identifier.UserID, userID)
}

func TestIdentifier_CreateFailure(t *testing.T) {
	// Create User
	userID := integration.CreateUser(t)
	// List identifiers
	identifier, err := integration.SDK(t).Identifiers().Create(context.Background(), userID, api.IdentifierCreateReq{IdentifierValue: "", IdentifierType: "email", Status: "verified"})
	assert.Error(t, err)
	assert.Nil(t, identifier)
}
