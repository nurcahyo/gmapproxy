package main

import (
	"errors"
	"log"
	"net"
	"sort"
	"strings"

	place "github.com/aniqma/gmapproxy/place"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) NearbySearchByTypes(ctx context.Context, in *place.Request) (*place.Response, error) {
	if in.GetCity() == "" || in.GetLatlong() == "" {
		return nil, errors.New("city, country and lat long are required")
	}
	if len(in.GetTypes()) < 1 {
		return nil, errors.New("Types can't be empty")
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
