package centralog

import (
	"context"
	"fmt"
	pb "github.com/XiovV/centralog-agent/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

func (a *App) ListNodesCmd() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTAINERS\tSTATUS\t")

	nodes, err := a.repository.GetNodes()
	if err != nil {
		log.Fatalln(err)
	}

	for _, node := range nodes {
		client := newClient(node.Location)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		response, err := client.GetRunningContainers(ctx, &pb.Containers{Containers: strings.Split(node.Containers, ",")})
		if err != nil {
			out := fmt.Sprintf("%s\t%d/%d\t%s", node.Name, 0, len(strings.Split(node.Containers, ",")), "DOWN")
			fmt.Fprintln(w, out)
			log.Fatalln(err)
		}

		out := fmt.Sprintf("%s\t%d/%d\t%s", node.Name, len(response.GetContainers()), len(strings.Split(node.Containers, ",")), "UP")
		fmt.Fprintln(w, out)
	}

	w.Flush()
}

func (a *App) ListContainersCmd(nodeName string) {
	node, err := a.repository.GetNode(nodeName)
	if err != nil {
		fmt.Println("node does not exist")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := a.centralogClient.GetContainersInfo(ctx, &pb.Containers{Containers: strings.Split(node.Containers, ",")})
	if err != nil {
		log.Fatalln(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")

	for _, container := range response.GetContainers() {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", container.Name, container.State))
	}

	w.Flush()
}

func (a *App) getContainersRPC(client pb.CentralogClient) []*pb.Container {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, _ := client.GetContainers(ctx, &pb.GetContainersRequest{})

	return response.Containers
}

func (a *App) pingServer(target string) (pb.CentralogClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewCentralogClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = client.Health(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (a *App) checkAPIKey(client pb.CentralogClient, key string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := client.CheckAPIKey(ctx, &pb.CheckAPIKeyRequest{Key: key})
	if err != nil {
		return err
	}

	return nil
}
