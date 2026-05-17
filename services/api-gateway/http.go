package main

import (
	"encoding/json"
	"net/http"
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
	// CHANGE THE LAST ARG TO "FALSE" if the OSRM API is not working right now
	route, err := tripService.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination, true)
	if err != nil {
		http.Error(w, "failed to get route from trip service", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: route}

	writeJSON(w, http.StatusOK, response)
}
