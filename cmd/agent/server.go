package main

import (
	"context"
	"fmt"
	docker2 "github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
	KEY_LENGTH     = 41
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.authenticateUnary), grpc.StreamInterceptor(s.authenticateStream))
	pb.RegisterCentralogServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

func (s *Server) Health(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: 1}, nil
}

func (s *Server) CheckAPIKey(ctx context.Context, in *pb.CheckAPIKeyRequest) (*pb.CheckAPIKeyResponse, error) {
	if len(in.GetKey()) > KEY_LENGTH {
		return &pb.CheckAPIKeyResponse{Valid: false}, status.Error(codes.InvalidArgument, "api key is too long")
	}

	key := s.Repository.GetAPIKey()

	err := bcrypt.CompareHashAndPassword(key, []byte(in.GetKey()))
	if err != nil {
		return &pb.CheckAPIKeyResponse{Valid: false}, status.Error(codes.InvalidArgument, "api key is invalid")
	}

	return &pb.CheckAPIKeyResponse{Valid: true}, nil
}
