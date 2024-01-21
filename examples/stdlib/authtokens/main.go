package main

import (
	"fmt"
	"net/http"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/generated/common"
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

	http.HandleFunc("/validateAuthToken", func(w http.ResponseWriter, r *http.Request) {
		corbadoAuthToken := r.URL.Query().Get("corbadoAuthToken")

		request := api.AuthTokenValidateReq{
			Token: corbadoAuthToken,
			ClientInfo: common.ClientInfo{
				RemoteAddress: "127.0.0.1",
				UserAgent:     "Corbado Go SDK Example",
			},
		}

		response, err := sdk.AuthTokens().Validate(r.Context(), request)
		if err != nil {
			if serverErr := corbado.AsServerError(err); serverErr != nil {
				fmt.Fprintf(w, serverErr.Error())
				return
			} else {
				panic(err)
			}
		}

		fmt.Fprintf(w, response.Data.UserID)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
