package main

//////////////////////////////////////////////////////////////////////////////////////////////
// Basic example which serves as basis for code snippets for integration guides             //
//////////////////////////////////////////////////////////////////////////////////////////////

import (
	"context"
	"fmt"
	"github.com/corbado/corbado-go"
	"net/http"
)

func main() {
	//////////////////////////////////////////////////////////////////////////////////////////////
	// Instantiate SDK                                                                          //
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Configuration
	projectID := "<Your Project ID>"
	apiSecret := "<Your API secret>"
	frontendApi := "<Your Frontend API URL>"
	backendApi := "<Your Backend API URL>"

	config, err := corbado.NewConfig(projectID, apiSecret, frontendApi, backendApi)
	if err != nil {
		panic(err)
	}

	// Create SDK instance
	sdk, err := corbado.NewSDK(config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////////////////////////////////////////////
		// Protecting routes                                                                        //
		//////////////////////////////////////////////////////////////////////////////////////////////

		// Retrieve the short-term session value from the Cookie (e.g. from cookie or Auth Headers)
		shortSession := "<Your short-term session value>"
		user, err := sdk.Sessions().ValidateToken(shortSession)

		if err != nil || user == nil {
			// User is not authenticated, redirect to login
			// page for example
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		// User is authenticated
		fmt.Fprintln(w, "User is authenticated!")

		//////////////////////////////////////////////////////////////////////////////////////////////
		// Getting user data from short-term session (represented as JWT)                           //
		//////////////////////////////////////////////////////////////////////////////////////////////

		usr, err := sdk.Sessions().ValidateToken(shortSession)

		if err != nil || usr == nil {
			// User is not authenticated, redirect to login
			// page for example
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		fmt.Fprintf(w, "User ID: %s\n", usr.UserID)
		fmt.Fprintf(w, "User full name: %s\n", usr.FullName)

		//////////////////////////////////////////////////////////////////////////////////////////////
		// Getting user data from Corbado Backend API                                               //
		//////////////////////////////////////////////////////////////////////////////////////////////

		newUser, err := sdk.Sessions().ValidateToken(shortSession)

		if err != nil || newUser == nil {
			// User is not authenticated, redirect to login
			// page for example
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		fullUser, err := sdk.Users().Get(context.Background(), newUser.UserID)

		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "User ID: %s\n", fullUser.UserID)
		fmt.Fprintf(w, "User full name: %s\n", *fullUser.FullName)
		fmt.Fprintf(w, "User status: %s\n", fullUser.Status)

		// To get the email we use the IdentifierService
		emailIdentifiers, err := sdk.Identifiers().ListByUserIDAndType(context.Background(), fullUser.UserID, "email", "", 1, 10)

		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "User Email: %s\n", emailIdentifiers.Identifiers[0].Value)

	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
