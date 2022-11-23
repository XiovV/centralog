package main

import (
	"errors"
	"github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/docker/docker/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

func (s *Server) GetLogs(request *pb.GetLogsRequest, stream pb.Centralog_GetLogsServer) error {
	err := s.validateGetLogsRequest(request)

	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	s.LogBuffer.Flush()

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true}

	if request.GetFollow() {
		wg := new(sync.WaitGroup)
		s.followLogs(wg, options, stream, request)

		wg.Wait()
		return nil
	}

	if request.GetLast() > 0 {
		logs, err := s.Repository.GetLastNLogs(request.GetLast())
		if err != nil {
			log.Fatalln(err)
		}

		for _, log := range logs {
			stream.Send(&pb.Log{
				Container: log.ContainerID,
				Timestamp: log.Timestamp,
				Message:   log.Message,
			})
		}
	}

	if request.GetFirst() > 0 {
		logs, err := s.Repository.GetFirstNLogs(request.GetFirst())
		if err != nil {
			log.Fatalln(err)
		}

		for _, log := range logs {
			stream.Send(&pb.Log{
				Container: log.ContainerID,
				Timestamp: log.Timestamp,
				Message:   log.Message,
			})
		}
	}

	return nil
}

func (s *Server) validateGetLogsRequest(request *pb.GetLogsRequest) error {
	if len(request.Containers) == 0 {
		return errors.New("containers array empty")
	}

	if request.GetFirst() > 0 && request.GetLast() > 0 {
		return errors.New("only first or last can be used, not both")
	}

	if (request.GetFirst() > 0 && request.GetShowAll()) || (request.GetLast() > 0 && request.GetShowAll()) {
		return errors.New("cannot use show all with first/last flags")
	}

	return nil
}

func (s *Server) followLogs(wg *sync.WaitGroup, options types.ContainerLogsOptions, stream pb.Centralog_GetLogsServer, request *pb.GetLogsRequest) {
	options.Follow = true

	if !request.ShowAll {
		options.Since = "0m"
	}

	stopSignals := []chan struct{}{}

	wg.Add(1)
	for _, container := range request.GetContainers() {
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
}
