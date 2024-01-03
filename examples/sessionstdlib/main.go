package main

import (
	"fmt"
	"net/http"

	"github.com/corbado/corbado-go"
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
		shortSession, err := sdkHelpers.GetShortSessionValue(r)
		if err != nil {
			panic(err)
		}

		user, err := sdk.Sessions().GetCurrentUser(shortSession)
		if err != nil {
			panic(err)
		}

		if user.Authenticated {
			// User is authenticated
			fmt.Fprint(w, "User is authenticated!")
		} else {
			// User is not authenticated, redirect to login
			// page for example
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
