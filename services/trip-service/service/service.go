package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrip(ctx context.Context, fare *domain.RideFareModal) (*domain.TripModal, error) {
	t := &domain.TripModal{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: *fare,
		Driver:   &pb.TripDriver{},
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *Service) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModal, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip fare: %w", err)
	}

	if fare == nil {
		return nil, fmt.Errorf("fare does not exist")
	}

	// User fare validation (user is owner of this fare?)
	if userID != fare.UserID {
		return nil, fmt.Errorf("fare does not belong to the user")
	}

	return fare, nil
}

func (s *Service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOsrmApi bool) (*tripTypes.OSRMResponse, error) {
	baseURL := "http://router.project-osrm.org"
	if !useOsrmApi {
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

	var osrmResp tripTypes.OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return nil, fmt.Errorf("failed to parse route response: %w", err)
	}
	return &osrmResp, nil
}

func (s *Service) EstimagePackagesPriceWithRoute(route *tripTypes.OSRMResponse) ([]*domain.RideFareModal, error) {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModal, len(baseFares))
	for i, fare := range baseFares {
		estimatedFares[i] = estimateFareRoute(fare, route)
	}
	return estimatedFares, nil
}

func (s *Service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModal, userID string) ([]*domain.RideFareModal, error) {
	fares := make([]*domain.RideFareModal, len(rideFares))
	for i, fare := range rideFares {
		id := primitive.NewObjectID()
		fares[i] = &domain.RideFareModal{
			ID:                id,
			UserID:            userID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}
		if err := s.repo.SaveRideFare(ctx, fares[i]); err != nil {
			return nil, err
		}
	}
	return fares, nil
}

func estimateFareRoute(f *domain.RideFareModal, route *tripTypes.OSRMResponse) *domain.RideFareModal {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKm * float64(pricingCfg.PricePerUnitDistance)
	timeFare := durationInMinutes * float64(pricingCfg.PricingPerUnitTime)
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModal{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModal {
	return []*domain.RideFareModal{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 500,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 600,
		},
		{
			PackageSlug:       "electric",
			TotalPriceInCents: 550,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 700,
		},
	}
}
