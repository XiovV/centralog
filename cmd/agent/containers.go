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
