package main

//////////////////////////////////////////////////////////////////////////////////////////////
// Basic example which serves as basis for code snippets for integration guides             //
//////////////////////////////////////////////////////////////////////////////////////////////

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/corbado/corbado-go"
)

func main() {
	//////////////////////////////////////////////////////////////////////////////////////////////
	// Instantiate SDK                                                                          //
	//////////////////////////////////////////////////////////////////////////////////////////////

	// Configuration
	projectID := "pro-8800191598924016720"
	apiSecret := "corbado1_UdPwavfp5AbWdFThUojMfTdY7Rm9Wu"
	frontendApi := "https://pro-8800191598924016720.frontendapi.cloud.corbado.io"
	backendApi := "https://backendapi.cloud.corbado.io"

	config, err := corbado.NewConfig(projectID, apiSecret, frontendApi, backendApi)
	if err != nil {
		panic(err)
	}

	// Create SDK instance
	sdk, err := corbado.NewSDK(config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/closedArea", func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////////////////////////////////////////////
		// Protecting routes                                                                        //
		//////////////////////////////////////////////////////////////////////////////////////////////

		// Retrieve short-term session value from cookie
		shortSessionCookie, err := r.Cookie("cbo_short_session")
		if errors.Is(err, http.ErrNoCookie) {
			// User is not authenticated, redirect to login page for example
			http.Redirect(w, r, "/login", http.StatusFound)

			return
		} else if err != nil {
			// Return full error (not recommended on production)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		shortSession := shortSessionCookie.Value

		user, err := sdk.Sessions().ValidateToken(shortSession)
		if err != nil {
			// Return full error (not recommended on production)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		// User is authenticated
		fmt.Fprintf(w, "User with ID %s is authenticated!", user.UserID)

		//////////////////////////////////////////////////////////////////////////////////////////////
		// Getting user data from short-term session (represented as JWT)                           //
		//////////////////////////////////////////////////////////////////////////////////////////////

		{
			user, err := sdk.Sessions().ValidateToken(shortSession)
			if err != nil {
				// Return full error (not recommended on production)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			fmt.Fprintf(w, "User ID: %s\n", user.UserID)
			fmt.Fprintf(w, "User full name: %s\n", user.FullName)
		}

		//////////////////////////////////////////////////////////////////////////////////////////////
		// Getting user data from Corbado Backend API                                               //
		//////////////////////////////////////////////////////////////////////////////////////////////

		{
			user, err := sdk.Sessions().ValidateToken(shortSession)
			if err != nil {
				// Return full error (not recommended on production)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			fullUser, err := sdk.Users().Get(context.Background(), user.UserID)
			if err != nil {
				// Return full error (not recommended on production)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			fmt.Fprintf(w, "User ID: %s\n", fullUser.UserID)

			if fullUser.FullName != nil {
				fmt.Fprintf(w, "User full name: %s\n", fullUser.FullName)
			}

			fmt.Fprintf(w, "User status: %s\n", fullUser.Status)

			// To get the email we use the identifier service
			emailIdentifiers, err := sdk.Identifiers().ListByUserIDAndType(context.Background(), fullUser.UserID, "email", "", 1, 10)
			if err != nil {
				// Return full error (not recommended on production)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			fmt.Fprintf(w, "User Email: %s\n", emailIdentifiers.Identifiers[0].Value)
		}
	})

	fmt.Println("Listening on :8000 ...")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
