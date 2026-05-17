package domain

import (
	"context"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModal struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare RideFareModal
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModal) (*TripModal, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModal) (*TripModal, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOsrmApi bool) (*types.OSRMResponse, error)
}
