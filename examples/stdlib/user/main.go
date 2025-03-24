package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/corbado/corbado-go/v2/pkg/generated/api"
	"log"
	"net/http"
	"os"

	"github.com/corbado/corbado-go/v2"
	"github.com/gorilla/mux"
)

var sdk corbado.SDK

// Create a new user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req api.UserCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := sdk.Users().Create(context.Background(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Create a new active user by full name
func createActiveByNameHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FullName string `json:"fullName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.FullName == "" {
		http.Error(w, "\"fullName\" is required", http.StatusBadRequest)
		return
	}

	user, err := sdk.Users().CreateActiveByName(context.Background(), body.FullName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating active user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Delete a user by ID
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if _, err := sdk.Users().Delete(context.Background(), userID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

// Get a user by ID
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := sdk.Users().Get(context.Background(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
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

	r.HandleFunc("/create", createUserHandler).Methods("POST")
	r.HandleFunc("/createActiveByName", createActiveByNameHandler).Methods("POST")
	r.HandleFunc("/delete/{id}", deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/get/{id}", getUserHandler).Methods("GET")

	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
