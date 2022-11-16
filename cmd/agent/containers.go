package main

import (
	"context"
	pb "github.com/XiovV/centralog-agent/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetContainers(ctx context.Context, in *pb.GetContainersRequest) (*pb.ContainerResponse, error) {
	containers, err := s.Docker.GetContainers()
	if err != nil {
		s.Logger.Error("couldn't get containers", zap.Error(err))
		return nil, status.Error(codes.Internal, "couldn't get containers")
	}

	response := []*pb.Container{}
	for _, container := range containers {
		response = append(response, &pb.Container{
			Id:    container.ID,
			Name:  container.Names[0][1:],
			State: container.State,
		})
	}

	return &pb.ContainerResponse{Containers: response}, nil
}

func (s *Server) GetRunningContainers(ctx context.Context, in *pb.Containers) (*pb.Containers, error) {
	runningContainers := &pb.Containers{}

	for _, container := range in.GetContainers() {
		c, err := s.Docker.GetContainer(container)
		if err == nil && c.State == "running" {
			runningContainers.Containers = append(runningContainers.Containers, container)
		}

	}

	return runningContainers, nil
}

func (s *Server) GetContainersInfo(ctx context.Context, in *pb.Containers) (*pb.ContainerResponse, error) {
	response := &pb.ContainerResponse{}

	for _, container := range in.GetContainers() {
		c, _ := s.Docker.GetContainer(container)

		response.Containers = append(response.Containers, &pb.Container{
			Id:    c.ID,
			Name:  c.Names[0][1:],
			State: c.State,
		})
	}

	return response, nil
}
