package main

import (
	"errors"
	"log"
	"net"
	"sort"
	"strings"

	place "github.com/nurcahyo/gmapproxy/place"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) NearbySearchByTypes(ctx context.Context, in *place.Request) (*place.Response, error) {
	if in.GetCountry() == "" || in.GetCity() == "" || in.GetLatlong() == "" {
		return nil, errors.New("city, country and lat long are required")
	}
	
	if in.GetTypes() == "" {
		return nil, errors.New("Types can't empty")
	}
	
	if in.GetKey() == "" {
		return nil, errors.New("API Key can't empty")
	}
	
	
	
	types := strings.Split(in.GetTypes(), ",")
	sort.Sort(sort.StringSlice(types))
	in.Types = strings.Join(types, ",")
	return place.FindNearbyPlaceByCityAndLatLong(in)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	place.RegisterPlaceServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
