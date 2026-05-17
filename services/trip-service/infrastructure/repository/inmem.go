package repository

import (
	"context"
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

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModal) (*domain.TripModal, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
