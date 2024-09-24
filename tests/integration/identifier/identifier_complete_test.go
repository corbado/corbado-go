//go:build integration

package identifier

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/tests/integration"
)

func TestIdentifierOperations(t *testing.T) {
	ctx := context.TODO()
	userID := integration.CreateUser(t)
	email := integration.CreateRandomTestEmail(t)
	initialIdentifier := &api.Identifier{}

	t.Run("CreateIdentifier", func(t *testing.T) {
		t.Run("ValidationError", func(t *testing.T) {
			identifier, err := integration.SDK(t).Identifiers().Create(ctx, userID, api.IdentifierCreateReq{
				IdentifierValue: "",
				IdentifierType:  "email",
				Status:          "verified",
			})
			assert.Error(t, err)
			assert.Nil(t, identifier)
		})

		t.Run("Success", func(t *testing.T) {
			identifier, err := integration.SDK(t).Identifiers().Create(ctx, userID, api.IdentifierCreateReq{
				IdentifierValue: email,
				IdentifierType:  "email",
				Status:          "verified",
			})
			assert.NoError(t, err)
			assert.NotNil(t, identifier)
			assert.Equal(t, identifier.Value, email)
			assert.Equal(t, identifier.Status, api.IdentifierStatus("verified"))
			assert.Equal(t, identifier.UserID, userID)

			initialIdentifier = identifier
		})
	})

	t.Run("ListIdentifiers", func(t *testing.T) {
		t.Run("All", func(t *testing.T) {

			// Create new identifier
			newIdentifier := integration.CreateIdentifier(t)

			var allIdentifiers []api.Identifier
			page := 1
			pageSize := 100
			hasNextPage := true

			for hasNextPage {
				// List identifiers for the current page
				list, err := integration.SDK(t).Identifiers().List(ctx, []string{}, "", page, pageSize)
				assert.NoError(t, err)
				assert.NotNil(t, list)

				// Append identifiers to the complete list
				allIdentifiers = append(allIdentifiers, list.Identifiers...)

				// Check if there's a next page
				if list.Paging.TotalPages <= page {
					hasNextPage = false
				} else {
					page++
				}
			}

			// Search for the new identifier in the list
			found := false
			for _, identifier := range allIdentifiers {
				if identifier.IdentifierID == newIdentifier {
					found = true
					break
				}
			}

			assert.True(t, found, "New identifier should be found in the list")
		})

		t.Run("ByValueAndType", func(t *testing.T) {
			list, err := integration.SDK(t).Identifiers().ListByValueAndType(ctx, email, "email", "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})

		t.Run("ByUserIDAndType", func(t *testing.T) {
			list, err := integration.SDK(t).Identifiers().ListByUserIDAndType(ctx, userID, "email", "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})

		t.Run("ByUserID", func(t *testing.T) {
			list, err := integration.SDK(t).Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.NotNil(t, list)
			assert.Len(t, list.Identifiers, 1)
			assert.Equal(t, list.Identifiers[0].Value, email)
		})
	})

	t.Run("UpdateIdentifier", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Update identifier status
			identifier, err := integration.SDK(t).Identifiers().UpdateStatus(ctx, userID, initialIdentifier.IdentifierID, api.IdentifierStatusPending)
			assert.NoError(t, err)
			assert.NotNil(t, identifier)

			// Verify the updated status
			list, err := integration.SDK(t).Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.Equal(t, list.Identifiers[0].Status, api.IdentifierStatusPending)
		})
	})

	t.Run("DeleteIdentifier", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// List identifiers before deletion
			initialList, err := integration.SDK(t).Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)

			// Delete identifier
			_, err = integration.SDK(t).Identifiers().Delete(ctx, userID, initialIdentifier.IdentifierID)
			assert.NoError(t, err)

			// List identifiers after deletion
			finalList, err := integration.SDK(t).Identifiers().ListByUserID(ctx, userID, "", 1, 100)
			assert.NoError(t, err)
			assert.Equal(t, len(finalList.Identifiers), len(initialList.Identifiers)-1)
		})

		t.Run("NotFound", func(t *testing.T) {
			// Attempt to delete a non-existent identifier
			identifier, err := integration.SDK(t).Identifiers().Delete(ctx, userID, initialIdentifier.IdentifierID)
			assert.Error(t, err)
			assert.Nil(t, identifier)
		})
	})
}
