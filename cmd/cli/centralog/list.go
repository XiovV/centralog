package centralog

import (
	"context"
	"fmt"
	pb "github.com/XiovV/centralog-agent/grpc"
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
		log.Fatalln("couldn't get nodes:", err)
	}

	for _, node := range nodes {
		err := a.initClient(node.Location)
		if err != nil {
			log.Fatalln("couldn't init client:", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		response, err := a.centralogClient.GetRunningContainers(ctx, &pb.Containers{Containers: strings.Split(node.Containers, ",")})
		if err != nil {
			out := fmt.Sprintf("%s\t%d/%d\t%s", node.Name, 0, len(strings.Split(node.Containers, ",")), "DOWN")
			fmt.Fprintln(w, out)
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

	a.initClient(node.Location)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := a.centralogClient.GetContainersInfo(ctx, &pb.Containers{Containers: strings.Split(node.Containers, ",")})
	if err != nil {
		log.Fatalln("couldn't get container info:", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")

	for _, container := range response.GetContainers() {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", container.Name, container.State))
	}

	w.Flush()
}
