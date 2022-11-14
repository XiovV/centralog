package main

import (
	"context"
	pb "github.com/XiovV/centralog-agent/grpc"
	"log"
)

func (s *Server) GetContainers(ctx context.Context, in *pb.GetContainersRequest) (*pb.ContainerResponse, error) {
	containers, err := s.Docker.GetContainers()
	if err != nil {
		log.Fatalln(err)
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

func (s *Server) GetRunningContainers(ctx context.Context, in *pb.RunningContainers) (*pb.RunningContainers, error) {
	runningContainers := &pb.RunningContainers{}

	for _, container := range in.GetContainers() {
		c, err := s.Docker.GetContainer(container)
		if err == nil && c.State == "running" {
			runningContainers.Containers = append(runningContainers.Containers, container)
		}

	}

	return runningContainers, nil
}

func (s *Server) GetContainersInfo(ctx context.Context, in *pb.RunningContainers) (*pb.ContainerResponse, error) {
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
