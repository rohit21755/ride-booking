package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/domain"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModal
	rideFares map[string]*domain.RideFareModal
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModal),
		rideFares: make(map[string]*domain.RideFareModal),
	}
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModal, error) {
	fare, exist := r.rideFares[id]
	if !exist {
		return nil, fmt.Errorf("fare does not exist with ID: %s", id)
	}

	return fare, nil
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModal) (*domain.TripModal, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModal) error {
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}
