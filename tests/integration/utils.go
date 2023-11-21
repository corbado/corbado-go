//go:build integration

package integration

import (
	"crypto/rand"
	"math/big"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/corbado/corbado-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func SDK(t *testing.T) corbado.SDK {
	config, err := corbado.NewConfig(GetProjectID(t), GetAPISecret(t))
	require.NoError(t, err)
	config.BackendAPI = GetBackendAPI(t)

	sdk, err := corbado.NewSDK(config)
	require.NoError(t, err)

	return sdk
}

func getEnv(t *testing.T, name string) string {
	env := os.Getenv(name)
	if env == "" {
		t.Fatalf("Missing env variable %s", name)
	}

	return env
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

func CreateRandomTestEmail(t *testing.T) string {
	value, err := generateString(10)
	require.NoError(t, err)

	return getFunctionName() + value + "@test.de"
}

func CreateRandomTestName(t *testing.T) string {
	value, err := generateString(10)
	require.NoError(t, err)

	return value
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

func getFunctionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)

	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	functionName := frame.Function
	functionName = functionName[strings.LastIndex(functionName, ".")+1:]
	functionName = functionName[5:]

	return functionName
}
