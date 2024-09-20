//go:build integration

package identifier

import (
	"context"
	"testing"

	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentifierOperations(t *testing.T) {
	sdkClient := integration.SDK(t)
	ctx := context.TODO()

	userID := integration.CreateUser(t)
	require.NotEmpty(t, userID)

	email := integration.CreateRandomTestEmail(t)

	initialIdentifier := api.Identifier{}

	// Subtest for Identifier Creation functionality
	t.Run("CreateIdentifier", func(t *testing.T) {
		t.Run("CreateSuccess", func(t *testing.T) {
			identifier, err := sdkClient.Identifiers().Create(ctx, userID, api.IdentifierCreateReq{
				IdentifierValue: email,
				IdentifierType:  "email",
				Status:          "verified",
			})
			assert.NoError(t, err)
			assert.NotNil(t, identifier)
			assert.Equal(t, identifier.Value, email)
			assert.Equal(t, identifier.Status, api.IdentifierStatus("verified"))
			assert.Equal(t, identifier.UserID, userID)

			initialIdentifier = *identifier
		})

		t.Run("CreateFailure", func(t *testing.T) {
			identifier, err := sdkClient.Identifiers().Create(ctx, userID, api.IdentifierCreateReq{
				IdentifierValue: "",
				IdentifierType:  "email",
				Status:          "verified",
			})
			assert.Error(t, err)
			assert.Nil(t, identifier)
		})
	})

	// Subtest for Identifier Listing functionality
	t.Run("ListIdentifiers", func(t *testing.T) {
		t.Run("ListAll", func(t *testing.T) {
			// List identifiers
			initialList, err := sdkClient.Identifiers().List(ctx, []string{}, "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, initialList)

			// Create new identifier
			integration.CreateIdentifier(t)

			// List identifiers containing new identifier
			newList, err := sdkClient.Identifiers().List(ctx, []string{}, "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, newList)
			assert.Equal(t, newList.Paging.TotalItems, initialList.Paging.TotalItems+1)
		})

		t.Run("ListByValueAndType", func(t *testing.T) {
			list, err := sdkClient.Identifiers().ListByValueAndType(ctx, email, "email", "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})

		t.Run("ListByUserIDAndType", func(t *testing.T) {
			list, err := sdkClient.Identifiers().ListByUserIDAndType(ctx, userID, "email", "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})

		t.Run("ListByUserID", func(t *testing.T) {
			list, err := sdkClient.Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})
	})

	// Subtest for Identifier Update functionality
	t.Run("UpdateIdentifier", func(t *testing.T) {
		t.Run("UpdateStatusSuccess", func(t *testing.T) {
			// Update identifier status
			identifier, err := sdkClient.Identifiers().UpdateStatus(ctx, userID, initialIdentifier.IdentifierID, api.IdentifierStatusPending)
			assert.NoError(t, err)
			assert.NotNil(t, identifier)

			// Verify the updated status
			list, err := sdkClient.Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.Equal(t, list.Identifiers[0].Status, api.IdentifierStatusPending)
		})
	})

	// Subtest for Identifier Deletion functionality
	t.Run("DeleteIdentifier", func(t *testing.T) {
		t.Run("DeleteSuccess", func(t *testing.T) {
			// List identifiers before deletion
			initialList, err := sdkClient.Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)

			// Delete identifier
			_, err = sdkClient.Identifiers().Delete(ctx, userID, initialIdentifier.IdentifierID)
			assert.NoError(t, err)

			// List identifiers after deletion
			finalList, err := sdkClient.Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.Equal(t, len(finalList.Identifiers), len(initialList.Identifiers)-1)
		})

		t.Run("DeleteFailure", func(t *testing.T) {
			// Attempt to delete a non-existent identifier
			identifier, err := sdkClient.Identifiers().Delete(ctx, userID, initialIdentifier.IdentifierID)
			assert.Error(t, err)
			assert.Nil(t, identifier)
		})
	})
}
