package main

import (
	"github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/docker/docker/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

func (s *Server) GetLogs(request *pb.GetLogsRequest, stream pb.Centralog_GetLogsServer) error {
	if len(request.Containers) == 0 {
		return status.Error(codes.InvalidArgument, "containers array empty")
	}

	s.LogBuffer.Flush()

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true}

	if request.GetFollow() {
		options.Follow = true

		if !request.ShowAll {
			options.Since = "0m"
		}

		stopSignals := []chan struct{}{}

		var wg sync.WaitGroup

		wg.Add(1)
		for _, container := range request.Containers {
			stopSignal := make(chan struct{})
			stopSignals = append(stopSignals, stopSignal)

			logWriter := docker.NewServerLogWriter(stopSignal, container, stream)

			go s.Docker.CollectLogs(container, logWriter, options)
		}

		// since the logs are being collected inside goroutines, we have to wait until
		// the client closes the connection so the connection doesn't get instantly
		// dropped by the server.
		select {
		case <-stream.Context().Done():
			for _, signal := range stopSignals {
				signal <- struct{}{}
			}
			wg.Done()
		}

		wg.Wait()
	}

	return nil
}
