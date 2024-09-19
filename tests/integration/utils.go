//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"math/big"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
)

func SDK(t *testing.T) corbado.SDK {
	config, err := corbado.NewConfig(GetProjectID(t), GetAPISecret(t), GetFrontendAPI(t), GetBackendAPI(t))
	require.NoError(t, err)

	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	return sdk
}

func GetProjectID(t *testing.T) string {
	return getEnv(t, "CORBADO_PROJECT_ID")
}

func GetAPISecret(t *testing.T) string {
	return getEnv(t, "CORBADO_API_SECRET")
}

func GetBackendAPI(t *testing.T) string {
	return getEnv(t, "CORBADO_BACKEND_API")
}

func GetFrontendAPI(t *testing.T) string {
	return getEnv(t, "CORBADO_FRONTEND_API")
}

func CreateRandomTestName(t *testing.T) *string {
	value, err := generateString(10)
	require.NoError(t, err)

	return &value
}

func CreateUser(t *testing.T) string {
	rsp, err := SDK(t).Users().Create(context.TODO(), api.UserCreateReq{
		FullName: CreateRandomTestName(t),
		Status:   "active",
	})
	require.NoError(t, err)

	return rsp.UserID
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

func generateNumber(length int) (string, error) {
	const letters = "0123456789"

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

func getEnv(t *testing.T, name string) string {
	env := os.Getenv(name)
	if env == "" {
		t.Fatalf("Missing env variable %s", name)
	}

	return env
}
