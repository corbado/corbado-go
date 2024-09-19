package main

//////////////////////////////////////////////////////////////////////////////////////////////
// Basic example which serves as basis for code snippets for integration guides             //
//////////////////////////////////////////////////////////////////////////////////////////////

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/corbado/corbado-go"
)

func main() {
	//////////////////////////////////////////////////////////////////////////////////////////////
	// Instantiate SDK                                                                          //
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Configuration
	projectID := os.Getenv("CORBADO_PROJECT_ID")
	apiSecret := os.Getenv("CORBADO_API_SECRET")
	frontendApi := os.Getenv("CORBADO_FRONTEND_API")
	backendApi := os.Getenv("CORBADO_BACKEND_API")

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
		shortSession := "eyJhbGciOiJSUzI1NiIsImtpZCI6InBraS01MDYxMTQ2NzY0ODg3ODc2OTMyIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwczovL2F1dGguY29yYmFkby1kZXYuY29tIiwic3ViIjoidXNyLTE1IiwiZXhwIjoxNzI2NzUwOTg2LCJuYmYiOjE3MjY3NTA2NzYsImlhdCI6MTcyNjc1MDY4NiwianRpIjoiUXVKZENrQ2x0R3QybVBGNlFlUmtTbmxZeEdLTXlzIiwibmFtZSI6IkFtaW5lIEhhbWRvdW5pIiwib3JpZyI6Im1vaGFtZWQuYW1pbmUuaGFtZG91bmlAY29yYmFkby5jb20iLCJlbWFpbCI6Im1vaGFtZWQuYW1pbmUuaGFtZG91bmlAY29yYmFkby5jb20iLCJ2ZXJzaW9uIjoyfQ.l3088ytJ-8LBmlzqqBH1n-ebY2yx56vfaQlRU5_eR7EwcMzAuwFIInsinTVlNlXOhf4s1l0YqMKMSkIitlk9c7eo09wV55ZZ76QXJh3NK3itkoNAkl8eaiszigbJecIExZuTzu7yG4l0gD9sq0Ik9eOD6pHN8WOLkImCkEGjORi-30HbS0oY8Kq4tpc3TJav4GIC9_PbVh075M97oyRn2Qza1q1PVwE5Xhh8jr01qn6tJynowhvWO1nBKnGtI0x8qBlapHx7jt7fmVWbVSHkCxJgnqUaCNZC9V-cdssOUhFK6BgP9JRYAL8uCEwOhd-NbLFoV-2R6VMx05KxvjlhnA"
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
