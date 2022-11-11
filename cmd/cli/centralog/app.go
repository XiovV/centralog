package centralog

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/XiovV/centralog-agent/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"reflect"
	"text/tabwriter"
)

type App struct {
	centralogClient pb.CentralogClient
}

func NewApp() *App {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewCentralogClient(conn)

	return &App{
		centralogClient: client,
	}
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

func (a *App) AddNodeWithPrompt() {
	var client pb.CentralogClient

	qs := []*survey.Question{
		{
			Name:   "url",
			Prompt: &survey.Input{Message: "Enter your node's URL:"},
			Validate: func(ans interface{}) error {
				val := reflect.ValueOf(ans).String()
				c, err := a.pingServer(val)
				if err != nil {
					return errors.New("connection refused, please check your URL.")
				}

				client = c

				return nil
			},
		},
		{
			Name:   "key",
			Prompt: &survey.Input{Message: "Enter your node's API key:"},
			Validate: func(ans interface{}) error {
				val := reflect.ValueOf(ans).String()

				err := a.checkAPIKey(client, val)
				if err != nil {
					return errors.New("api key is invalid")
				}

				return nil
			},
		},
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Enter your node's custom name:"},
			Validate: survey.Required,
		},
	}

	var answers struct {
		Url  string
		Key  string
		Name string
	}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalln(err)
	}

	containers := []string{}
	for _, container := range a.getContainersRPC(client) {
		containers = append(containers, fmt.Sprintf("%s (%s)", container.Name, container.State))
	}

	prompt := &survey.MultiSelect{
		Message: "Select containers:",
		Options: containers,
	}

	survey.AskOne(prompt, &containers)

	fmt.Printf("Node %s added successfully\n", answers.Name)
	fmt.Println("You selected:", containers)
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}
