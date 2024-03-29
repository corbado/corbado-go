package main

import (
	"context"
	"fmt"

	"github.com/corbado/corbado-go"
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

	// list all users
	users, err := sdk.Users().List(context.TODO(), nil)
	if err != nil {
		// handle server errors and client errors differently
		if serverErr := corbado.AsServerError(err); serverErr != nil {
			fmt.Printf("Received server error: %s", serverErr)

			return
		} else {
			panic(err)
		}
	}

	for _, usr := range users.Data.Users {
		fmt.Printf("%s: %s\n", usr.ID, usr.Name)
	}
}
