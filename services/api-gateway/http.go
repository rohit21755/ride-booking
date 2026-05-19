package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/services/trip-service/domain"
	"ride-sharing/shared/contracts"
)

var tripService domain.TripService

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	// Call trip service
	ctx := r.Context()

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripServiceClient.Close()
	// CHANGE THE LAST ARG TO "FALSE" if the OSRM API is not working right now

	tripPreview, err := tripServiceClient.Client.PreviewTrip(ctx, reqBody.ToProto())
	if err != nil {
		http.Error(w, "failed to get route from trip service", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreview}

	writeJSON(w, http.StatusOK, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody createTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" || reqBody.RideFareID == "" {
		http.Error(w, "user ID and ride fare ID are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Printf("failed to connect to trip service: %v", err)
		http.Error(w, "failed to connect to trip service", http.StatusInternalServerError)
		return
	}
	defer tripServiceClient.Close()

	tripResponse, err := tripServiceClient.Client.CreateTrip(ctx, reqBody.ToProto())
	if err != nil {
		log.Printf("failed to create trip: %v", err)
		http.Error(w, "failed to create trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripResponse}
	writeJSON(w, http.StatusOK, response)
}
