package main

import (
	"fmt"
	"net/http"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/stdlib"
)

func main() {
	// NewConfigEnv() reads project ID and API secret from CORBADO_PROJECT_ID
	// and CORBADO_API_SECRET environment variables
	config, err := corbado.NewConfigEnv()
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
			fmt.Fprintf(w, "Hello %s!", user.Name)
		} else {
			fmt.Fprintf(w, "Hello guest!")
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
