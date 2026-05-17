package http

import (
	"encoding/json"
	"net/http"
	"ride-sharing/services/trip-service/domain"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service domain.TripService
}

type previewTripRequest struct {
	UserID      string           `json:"user_id"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (s *HttpHandler) HandlePreviewTrip(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// CHANGE THE LAST ARG TO "FALSE" if the OSRM API is not working right now
	t, err := s.Service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination, true)
	if err != nil {
		http.Error(w, "Failed to create trip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
