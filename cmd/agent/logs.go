package main

import (
	docker2 "github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/docker/docker/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

func (s *Server) FollowLogs(request *pb.FollowLogsRequest, stream pb.Centralog_FollowLogsServer) error {
	if len(request.Containers) == 0 {
		return status.Error(codes.InvalidArgument, "containers array empty")
	}

	s.LogBuffer.Flush()

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
