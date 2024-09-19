package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/stdlib"
)

func main() {
	projectID := os.Getenv("CORBADO_PROJECT_ID")
	apiSecret := os.Getenv("CORBADO_API_SECRET")
	frontendApi := os.Getenv("CORBADO_FRONTEND_API")
	backendApi := os.Getenv("CORBADO_BACKEND_API")

	config, err := corbado.NewConfig(projectID, apiSecret, frontendApi, backendApi)
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

	http.HandleFunc("/logged-in", func(w http.ResponseWriter, r *http.Request) {
		shortSession, err := sdkHelpers.GetShortSessionValue(r, "cbo_short_session")
		if err != nil {
			panic(err)
		}

		user, err := sdk.Sessions().ValidateToken(shortSession)

		if err != nil || user == nil {
			// User is not authenticated, redirect to login
			// page for example

			http.Error(w, "User is not authenticated!", http.StatusUnauthorized)
			return
		}

		// User is authenticated
		fmt.Fprintf(w, "User is authenticated! \n")

		fmt.Fprintf(w, "User ID: %s\n", user.UserID)
		fmt.Fprintf(w, "User full name: %s\n", user.FullName)
	})

	http.HandleFunc("/set-cookie", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("session")

		http.SetCookie(w, &http.Cookie{
			Name:  "cbo_short_session",
			Value: value,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cookie set"))
	})

	if err := http.ListenAndServe(":8100", nil); err != nil {
		panic(err)
	}
}
