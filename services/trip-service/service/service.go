package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrip(ctx context.Context, fare *domain.RideFareModal) (*domain.TripModal, error) {
	return s.repo.CreateTrip(ctx, &domain.TripModal{
		ID:       primitive.NewObjectID(),
		UserID:   "",
		Status:   "",
		RideFare: *fare,
	})
}

func (s *Service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOsrmApi bool) (*types.OSRMResponse, error) {
	baseURL := "http://router.project-osrm.org"
	if !useOsrmApi {
		// In a real scenario, this might point to a self-hosted OSRM instance
		// or use a fallback mechanism. For now, we'll just keep using the base URL
		// but this structure allows for easy toggling.
		log.Println("OSRM API toggle is FALSE, using default fallback (if any)")
	}

	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", baseURL, pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read route response: %w", err)
	}

	var osrmResp types.OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return nil, fmt.Errorf("failed to parse route response: %w", err)
	}
	return &osrmResp, nil
}
