package grpc

import (
	"context"
	"ride-sharing/services/trip-service/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	tripSvc domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		tripSvc: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()

	rideFare, err := h.tripSvc.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, err
	}

	trip, err := h.tripSvc.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, err
	}

	// Add a comment at the end of the function to publish an event on the Async Comms module.

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
		Trip:   trip.ToProto(),
	}, nil
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStart()
	destination := req.GetDestination()

	pickupCorr := &types.Coordinate{
		Latitude:  pickup.GetLatitude(),
		Longitude: pickup.GetLongitude(),
	}

	destCorr := &types.Coordinate{
		Latitude:  destination.GetLatitude(),
		Longitude: destination.GetLongitude(),
	}

	route, err := h.tripSvc.GetRoute(ctx, pickupCorr, destCorr, true)
	if err != nil {
		return nil, err
	}

	estimatedFares, err := h.tripSvc.EstimagePackagesPriceWithRoute(route)
	if err != nil {
		return nil, err
	}
	fares, err := h.tripSvc.GenerateTripFares(ctx, estimatedFares, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}
