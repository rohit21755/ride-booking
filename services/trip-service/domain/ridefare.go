package domain

import (
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModal struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
}

func (r *RideFareModal) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(rideFares []*RideFareModal) []*pb.RideFare {
	protoFares := make([]*pb.RideFare, len(rideFares))
	for i, rideFare := range rideFares {
		protoFares[i] = rideFare.ToProto()
	}
	return protoFares
}
