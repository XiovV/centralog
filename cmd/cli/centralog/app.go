package centralog

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

type App struct {
	centralogClient pb.CentralogClient
	repository      *repository.SQLite
}

func NewApp() *App {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewCentralogClient(conn)

	repo := repository.New()

	return &App{
		centralogClient: client,
		repository:      repo,
	}
}

func newClient(target string) pb.CentralogClient {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewCentralogClient(conn)
}

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
		response, err := client.GetRunningContainers(ctx, &pb.RunningContainers{Containers: strings.Split(node.Containers, ",")})
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
	response, err := a.centralogClient.GetContainersInfo(ctx, &pb.RunningContainers{Containers: strings.Split(node.Containers, ",")})
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
			Name:   "name",
			Prompt: &survey.Input{Message: "Enter your node's custom name:"},
			Validate: func(ans interface{}) error {
				val := reflect.ValueOf(ans).String()

				exists := a.repository.DoesNodeExist(val)
				if exists {
					return errors.New("a node with this name already exists, please choose a different name.")
				}

				return nil
			},
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

	containersList := []string{}
	for _, container := range a.getContainersRPC(client) {
		containersList = append(containersList, fmt.Sprintf("%s (%s)", container.Name, container.State))
	}

	prompt := &survey.MultiSelect{
		Message: "Select containers:",
		Options: containersList,
	}

	containersSelected := []string{}
	survey.AskOne(prompt, &containersSelected)

	fmt.Printf("Node %s added successfully\n", answers.Name)

	containers := []string{}
	for _, container := range containersSelected {
		containers = append(containers, strings.Split(container, " ")[0])
	}

	node := repository.Node{
		Location:   answers.Url,
		APIKey:     answers.Key,
		Name:       answers.Name,
		Containers: strings.Join(containers, ","),
	}

	a.repository.InsertNode(node)
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}
