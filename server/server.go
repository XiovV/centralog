package server

import (
	"fmt"
	docker2 "github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/docker/docker/api/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"sync"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
)

type Server struct {
	pb.UnimplementedLogsServer

	Logger     *zap.Logger
	Docker     *docker2.Controller
	Repository *repository.Repository
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogsServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

func (s *Server) FollowLogs(request *pb.FollowLogsRequest, stream pb.Logs_FollowLogsServer) error {
	if len(request.Containers) == 0 {
		return status.Error(codes.InvalidArgument, "containers array empty")
	}

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true}

	if !request.ShowAll {
		options.Since = "0m"
	}

	for _, container := range request.Containers {
		logWriter := docker2.NewServerLogWriter(container, stream)

		go s.Docker.CollectLogs(container, logWriter, options)
	}

	var wg sync.WaitGroup

	// since the logs are being collected inside goroutines, we have to wait until
	// the client closes the connection so the connection doesn't get instantly
	// dropped by the server.
	wg.Add(1)
	select {
	case <-stream.Context().Done():
		wg.Done()
	}

	wg.Wait()

	return nil
}
