package types

import (
	pb "ride-sharing/shared/proto/trip"
)

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OSRMResponse) ToProto() *pb.Route {
	route := o.Routes[0]
	geometery := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometery))
	for i, coord := range geometery {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coord[1],
			Longitude: coord[0],
		}
	}

	return &pb.Route{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: []*pb.Geometry{
			&pb.Geometry{
				Coordinates: coordinates,
			},
		},
	}
}

type PricingConfig struct {
	PricePerUnitDistance float32
	PricingPerUnitTime   float32
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitDistance: 1.5,
		PricingPerUnitTime:   0.25,
	}
}
