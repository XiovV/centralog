package centralog

import (
	"context"
	"fmt"
	pb "github.com/XiovV/centralog-agent/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
)

func (a *App) ShowLogs(nodeName string, containersFlag []string) {
	node, err := a.repository.GetNode(nodeName)
	if err != nil {
		log.Fatalln(err)
	}

	err = a.initClient(node.Location)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", node.APIKey)

	containers := node.GetContainers()
	if len(containersFlag) >= 1 {
		containers = containersFlag
	}

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

		timestamp := time.UnixMilli(resp.Timestamp)

		fmt.Printf("container: %s | timestamp: %s | message: %s", resp.Container, timestamp.Format("02/01/2006, 15:04:05"), resp.Message)
	}
}
