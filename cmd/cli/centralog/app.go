package centralog

import (
	"context"
	"fmt"
	"github.com/XiovV/centralog-agent/cmd/cli/pkg/prompt"
	pb "github.com/XiovV/centralog-agent/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

type App struct {
	client pb.LogsClient
}

func NewApp() *App {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewLogsClient(conn)

	return &App{client: client}
}

func (a *App) ListNodesCmd() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTAINERS\tSTATUS\t")
	fmt.Fprintln(w, "node1\t2/3\tUP\t")
	fmt.Fprintln(w, "myNewNode\t0/5\tDOWN\t")

	w.Flush()
}

func (a *App) ListContainersCmd(node string) {
	if node != "node1" && node != "myNewNode" {
		fmt.Println(node, "does not exist")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")
	if node == "node1" {
		fmt.Fprintln(w, "instance1\tUP")
		fmt.Fprintln(w, "instance2\tUP")
		fmt.Fprintln(w, "instance3\tDOWN")
	} else {
		fmt.Fprintln(w, "instance1\tDOWN")
		fmt.Fprintln(w, "instance2\tDOWN")
		fmt.Fprintln(w, "instance3\tDOWN")
		fmt.Fprintln(w, "instance4\tDOWN")
		fmt.Fprintln(w, "instance5\tDOWN")
	}

	w.Flush()
}

func (a *App) ShowLogs() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &pb.FollowLogsRequest{
		Containers: []string{"logserver1", "logserver2"},
		ShowAll:    false,
	}

	stream, err := a.client.FollowLogs(ctx, request)
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

func (a *App) AddNodeWithPrompt() {
	prompt.GetURL("Enter your node's URL")
	_, err := prompt.GetAPIKey("Enter your node's API key")
	if err != nil {
		fmt.Println(err)
	}

	name := prompt.GetString("Enter your node's custom name")

	fmt.Printf("Node %s added successfully\n", name)
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}
