package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/stdlib"
)

func main() {
	// NewConfigFromEnv() reads project ID and API secret from CORBADO_PROJECT_ID
	// and CORBADO_API_SECRET environment variables
	config, err := corbado.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	sdk, err := corbado.NewSDK(config)
	if err != nil {
		panic(err)
	}

	sdkHelpers, err := stdlib.NewSDKHelpers(config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortSession, err := sdkHelpers.GetShortSessionValue(r, "cbo_short_session")
		if err != nil {
			panic(err)
		}

		user, err := sdk.Sessions().ValidateToken(shortSession)
		if err != nil {
			// User is not authenticated, redirect to login
			// page for example

			http.Redirect(w, r, "/login", http.StatusFound)
		}

		// User is authenticated
		fmt.Fprint(w, "User is authenticated!")

		fmt.Fprintf(w, "User ID: %s\n", user.UserID)
		fmt.Fprintf(w, "User full name: %s\n", user.FullName)

		rsp, err := sdk.Users().Get(context.Background(), user.UserID, nil)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "User status: %s\n", rsp.Status)

		identifiers, err := sdk.Identifiers().ListByUserIDAndType(context.Background(), user.UserID, api.IdentifierType("email"), "", 1, 10, nil)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "User email: %s\n", identifiers.Identifiers[0].Value)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
