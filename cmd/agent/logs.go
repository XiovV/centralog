package main

import (
	"errors"
	"github.com/XiovV/centralog-agent/docker"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/docker/docker/api/types"
	"go.uber.org/zap"
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

	if request.GetFirst() > 0 || request.GetLast() > 0 {
		logs, err := s.getNLogs(request.GetFirst(), request.GetLast())
		if err != nil {
			log.Fatalln(err)
		}

		s.sendLogs(logs, stream)

		if request.GetFollow() {
			wg := new(sync.WaitGroup)
			s.followLogs(wg, options, stream, request)

			wg.Wait()
		}

		return nil
	}

	if request.GetFollow() {
		wg := new(sync.WaitGroup)
		s.followLogs(wg, options, stream, request)

		wg.Wait()
		return nil
	}

	return nil
}

func (s *Server) getNLogs(first, last int32) ([]repository.Log, error) {
	if first > 0 {
		logs, err := s.Repository.GetFirstNLogs(first)
		if err != nil {
			return nil, err
		}

		return logs, nil
	}

	logs, err := s.Repository.GetLastNLogs(last)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *Server) sendLogs(logs []repository.Log, stream pb.Centralog_GetLogsServer) {
	for _, log := range logs {
		stream.Send(&pb.Log{
			Container: log.ContainerID,
			Timestamp: log.Timestamp,
			Message:   log.Message,
		})
	}
}

func (s *Server) validateGetLogsRequest(request *pb.GetLogsRequest) error {
	if request.GetFirst() > 0 && request.GetLast() > 0 {
		return errors.New("only first or last can be used, not both")
	}

	if (request.GetFirst() > 0 && request.GetShowAll()) || (request.GetLast() > 0 && request.GetShowAll()) {
		return errors.New("cannot use show all with first/last flags")
	}

	return nil
}

func (s *Server) followLogs(wg *sync.WaitGroup, options types.ContainerLogsOptions, stream pb.Centralog_GetLogsServer, request *pb.GetLogsRequest) {
	s.Logger.Info("follow logs")
	options.Follow = true

	if request.ShowAll {
		logs, err := s.Repository.GetAllLogs()
		if err != nil {
			log.Fatalln(err)
		}

		s.sendLogs(logs, stream)
	}

	options.Since = "0m"

	stopSignals := []chan struct{}{}

	wg.Add(1)

	containers := []string{}
	if len(request.GetContainers()) > 0 {
		containers = request.GetContainers()
	} else {
		config, err := s.Repository.GetConfig()
		if err != nil {
			s.Logger.Error("couldn't read config", zap.Error(err))
			return
		}

		containers = config.GetContainers()
	}

	for _, container := range containers {
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
