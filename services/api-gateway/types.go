package main

import (
	pb "ride-sharing/shared/proto/trip"
	types "ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (p *previewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserId:      p.UserID,
		Start:       &pb.Coordinate{Latitude: p.Pickup.Latitude, Longitude: p.Pickup.Longitude},
		Destination: &pb.Coordinate{Latitude: p.Destination.Latitude, Longitude: p.Destination.Longitude},
	}
}

type createTripRequest struct {
	UserID     string `json:"userID"`
	RideFareID string `json:"rideFareID"`
}

func (p *createTripRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		UserID:     p.UserID,
		RideFareID: p.RideFareID,
	}
}
