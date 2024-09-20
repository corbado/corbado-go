//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
)

func SDK(t *testing.T) corbado.SDK {
	config, err := corbado.NewConfigFromEnv()
	require.NoError(t, err)

	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	return sdk
}

func CreateRandomTestName(t *testing.T) *string {
	value, err := generateString(10)
	require.NoError(t, err)

	return &value
}

func CreateRandomTestEmail(t *testing.T) string {
	value, err := generateString(10)
	require.NoError(t, err)

	return "integration-test+" + value + "@corbado.com"
}

func CreateUser(t *testing.T) string {
	rsp, err := SDK(t).Users().Create(context.TODO(), api.UserCreateReq{
		FullName: CreateRandomTestName(t),
		Status:   "active",
	})
	require.NoError(t, err)

	return rsp.UserID
}

func CreateIdentifier(t *testing.T) (string, string, string) {
	userId := CreateUser(t)

	email := CreateRandomTestEmail(t)

	rsp, err := SDK(t).Identifiers().Create(context.TODO(), userId, api.IdentifierCreateReq{
		IdentifierType:  "email",
		IdentifierValue: email,
		Status:          "verified",
	})

	require.NoError(t, err)

	return rsp.IdentifierID, userId, email
}

func generateString(length int) (string, error) {
	// Removed I, 1, 0 and O because of risk of confusion
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopwrstuvwxyz23456789"

	res := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", errors.WithStack(err)
		}

		res[i] = letters[num.Int64()]
	}

	return string(res), nil
}
