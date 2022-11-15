package centralog

import (
	"context"
	"fmt"
	pb "github.com/XiovV/centralog-agent/grpc"
	"io"
	"log"
)

func (a *App) ShowLogs(containers []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &pb.FollowLogsRequest{
		Containers: containers,
		ShowAll:    false,
	}

	stream, err := a.centralogClient.FollowLogs(ctx, request)
	if err != nil {
		log.Fatalf("error initialising stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}

		fmt.Printf("container: %s | timestamp: %d | message: %s", resp.Container, resp.Timestamp, resp.Message)
	}
}
