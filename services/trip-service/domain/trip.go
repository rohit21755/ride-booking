package domain

import (
	"context"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	pb "ride-sharing/shared/proto/trip"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModal struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare RideFareModal
	Driver   *pb.TripDriver
}

func (t *TripModal) ToProto() *pb.Trip {
	return &pb.Trip{
		Id:           t.ID.Hex(),
		SelectedFare: t.RideFare.ToProto(),
		Status:       t.Status,
		UserID:       t.UserID,
		Driver:       t.Driver,
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModal) (*TripModal, error)
	SaveRideFare(ctx context.Context, fare *RideFareModal) error
	GetRideFareByID(ctx context.Context, rideFareID string) (*RideFareModal, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModal) (*TripModal, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOsrmApi bool) (*tripTypes.OSRMResponse, error)
	EstimagePackagesPriceWithRoute(route *tripTypes.OSRMResponse) ([]*RideFareModal, error)
	GenerateTripFares(ctx context.Context, rideFares []*RideFareModal, userID string) ([]*RideFareModal, error)

	GetAndValidateFare(ctx context.Context, rideFareID string, userID string) (*RideFareModal, error)
}
