package centralog

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"google.golang.org/grpc/metadata"
	"log"
	"reflect"
	"strings"
)

func (a *App) AddNodeWithPrompt() {
	qs := []*survey.Question{
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "Enter your node's URL:"},
			Validate: a.validateURL,
		},
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Enter your node's API key:"},
			Validate: a.validateKey,
		},
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Enter your node's custom name:"},
			Validate: a.validateNodeName,
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

	nodeContainers, err := a.getNodeContainers(answers.Key)
	if err != nil {
		log.Fatalln("couldn't fetch containers:", err)
	}

	containersList := []string{}
	for _, container := range nodeContainers {
		containersList = append(containersList, fmt.Sprintf("%s (%s)", container.Name, container.State))
	}

	prompt := &survey.MultiSelect{
		Message: "Select containers:",
		Options: containersList,
	}

	containersSelected := []string{}
	survey.AskOne(prompt, &containersSelected)

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

	err = a.repository.InsertNode(node)
	if err != nil {
		log.Fatalln("couldn't insert node:", err)
	}

	fmt.Printf("Node %s added successfully\n", answers.Name)
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}

func (a *App) validateURL(ans interface{}) error {
	val := reflect.ValueOf(ans).String()

	err := a.initClient(val)
	if err != nil {
		log.Fatalln("couldn't initiate client:", err)
	}

	err = a.pingServer()
	if err != nil {
		return errors.New("connection refused, please check your URL.")
	}

	return nil
}

func (a *App) validateKey(ans interface{}) error {
	val := reflect.ValueOf(ans).String()

	err := a.checkAPIKey(val)
	if err != nil {
		return errors.New("api key is invalid")
	}

	return nil
}

func (a *App) validateNodeName(ans interface{}) error {
	val := reflect.ValueOf(ans).String()

	if val == "" {
		return errors.New("please provide a node name")
	}

	if strings.Contains(val, " ") {
		return errors.New("node name should not contain spaces")
	}

	exists := a.repository.DoesNodeExist(val)
	if exists {
		return errors.New("a node with this name already exists, please choose a different name.")
	}

	return nil
}

func (a *App) getNodeContainers(apiKey string) ([]*pb.Container, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "sdfsfsdf")

	response, err := a.centralogClient.GetContainers(ctx, &pb.GetContainersRequest{})
	if err != nil {
		return nil, err
	}

	return response.Containers, nil
}

func (a *App) pingServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.centralogClient.Health(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return err
	}

	return nil
}

func (a *App) checkAPIKey(key string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.centralogClient.CheckAPIKey(ctx, &pb.CheckAPIKeyRequest{Key: key})
	if err != nil {
		return err
	}

	return nil
}
