package server

import (
	"context"
	"fmt"
	docker2 "github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
)

type Server struct {
	pb.UnimplementedCentralogServer

	Logger     *zap.Logger
	Docker     *docker2.Controller
	Repository *repository.SQLite
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCentralogServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

func (s *Server) Health(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: 1}, nil
}

func (s *Server) CheckAPIKey(ctx context.Context, in *pb.CheckAPIKeyRequest) (*pb.CheckAPIKeyResponse, error) {
	// TODO: hash the api key and compare the hash
	key := s.Repository.GetAPIKey()
	fmt.Println(key)

	return &pb.CheckAPIKeyResponse{Valid: false}, nil

	// TODO: return an appropriate error
	//if key != in.Key {
	//	return &pb.CheckAPIKeyResponse{Valid: false}, nil
	//}
	//
	//return &pb.CheckAPIKeyResponse{Valid: true}, nil
}
