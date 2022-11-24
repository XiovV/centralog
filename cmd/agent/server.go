package main

import (
	"context"
	"fmt"
	"github.com/XiovV/centralog-agent/docker"
	docker2 "github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/docker/docker/api/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"time"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
	KEY_LENGTH     = 41
	RATE_LIMTI     = 5
)

type Server struct {
	pb.UnimplementedCentralogServer

	Logger     *zap.Logger
	Docker     *docker2.Controller
	Repository *repository.SQLite
	LogBuffer  *docker.LogBuffer
	RateLimit  *Limiter

	apiKey []byte
}

func (s *Server) Serve() error {
	s.apiKey = s.Repository.GetAPIKey()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.authenticateUnary), grpc.StreamInterceptor(s.authenticateStream))
	pb.RegisterCentralogServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

func (s *Server) ListenForLogs() error {
	config, err := s.Repository.GetConfig()
	if err != nil {
		return err
	}

	containers := config.GetContainers()

	for _, containerName := range containers {
		options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true}

		s.Logger.Info("getting container", zap.String("containerName", containerName))

		_, err := s.Docker.GetContainer(containerName)

		if err != nil {
			s.Logger.Info("couldn't get container", zap.String("containerName", containerName), zap.Error(err))
			return nil
		}

		timestamp, err := s.Repository.GetLastTimestamp(containerName)
		if err != nil {
			s.Logger.Warn("couldn't get last timestamp:", zap.Error(err))
			return err
		}

		options.Since = time.UnixMilli(timestamp + 1).Format("2006-01-02T15:04:05.999999999")

		logWriter := docker.NewBackgroundLogWriter(s.LogBuffer, containerName)

		go s.Docker.CollectLogs(containerName, logWriter, options)
	}

	return nil
}

func (s *Server) Health(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: 1}, nil
}

func (s *Server) CheckAPIKey(ctx context.Context, in *pb.CheckAPIKeyRequest) (*pb.CheckAPIKeyResponse, error) {
	ok := s.RateLimit.Try()
	if !ok {
		s.Logger.Warn("check api key failed attempts limit reached")
		return &pb.CheckAPIKeyResponse{Valid: false}, status.Error(codes.InvalidArgument, fmt.Sprintf("rate limit exceeded, please try again in a few minutes"))
	}

	if len(in.GetKey()) > KEY_LENGTH {
		s.RateLimit.Fail()
		return &pb.CheckAPIKeyResponse{Valid: false}, status.Error(codes.InvalidArgument, "api key is too long")
	}

	err := bcrypt.CompareHashAndPassword(s.apiKey, []byte(in.GetKey()))
	if err != nil {
		s.RateLimit.Fail()
		s.Logger.Warn("api key is incorrect")
		return &pb.CheckAPIKeyResponse{Valid: false}, status.Error(codes.InvalidArgument, "api key is invalid")
	}

	return &pb.CheckAPIKeyResponse{Valid: true}, nil
}
