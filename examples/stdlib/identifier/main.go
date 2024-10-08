package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/corbado/corbado-go/v2"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"

	"github.com/gorilla/mux"
)

var sdk corbado.SDK

// Create a new identifier for a user
func createIdentifierHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var req api.IdentifierCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	identifier, err := sdk.Identifiers().Create(context.Background(), userID, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating identifier: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(identifier)
}

// Delete an identifier by user ID and identifier ID
func deleteIdentifierHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	identifierID := vars["identifierID"]

	result, err := sdk.Identifiers().Delete(context.Background(), userID, identifierID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting identifier: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// List identifiers with optional filters, sorting, pagination
func listIdentifiersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query["filter"]
	sort := query.Get("sort")
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))

	identifiers, err := sdk.Identifiers().List(context.Background(), filter, sort, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing identifiers: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identifiers)
}

// List identifiers by value and type
func listIdentifiersByValueAndTypeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	value := query.Get("value")
	identifierType := query.Get("type")
	sort := query.Get("sort")
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))

	if value == "" || identifierType == "" {
		http.Error(w, "\"value\" and \"type\" are required", http.StatusBadRequest)
		return
	}

	identifiers, err := sdk.Identifiers().ListByValueAndType(context.Background(), value, api.IdentifierType(identifierType), sort, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing identifiers by value and type: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identifiers)
}

// List identifiers by user ID
func listIdentifiersByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	query := r.URL.Query()
	sort := query.Get("sort")
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))

	identifiers, err := sdk.Identifiers().ListByUserID(context.Background(), userID, sort, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing identifiers by user ID: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identifiers)
}

// Update the status of an identifier
func updateIdentifierStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	identifierID := vars["identifierID"]

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		http.Error(w, "\"status\" is required", http.StatusBadRequest)
		return
	}

	identifier, err := sdk.Identifiers().UpdateStatus(context.Background(), userID, identifierID, api.IdentifierStatus(req.Status))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating identifier status: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identifier)
}

func main() {
	r := mux.NewRouter()

	projectID := os.Getenv("CORBADO_PROJECT_ID")
	apiSecret := os.Getenv("CORBADO_API_SECRET")
	frontendApi := os.Getenv("CORBADO_FRONTEND_API")
	backendApi := os.Getenv("CORBADO_BACKEND_API")

	config, err := corbado.NewConfig(projectID, apiSecret, frontendApi, backendApi)
	if err != nil {
		panic(err)
	}

	newSdk, err := corbado.NewSDK(config)
	if err != nil {
		panic(err)
	}

	sdk = newSdk

	r.HandleFunc("/create/{userID}", createIdentifierHandler).Methods("POST")
	r.HandleFunc("/delete/{userID}/{identifierID}", deleteIdentifierHandler).Methods("DELETE")
	r.HandleFunc("/list", listIdentifiersHandler).Methods("GET")
	r.HandleFunc("/listByValueAndType", listIdentifiersByValueAndTypeHandler).Methods("GET")
	r.HandleFunc("/listByUserId/{userID}", listIdentifiersByUserIDHandler).Methods("GET")
	r.HandleFunc("/updateStatus/{userID}/{identifierID}", updateIdentifierStatusHandler).Methods("PUT")

	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
