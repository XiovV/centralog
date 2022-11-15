package centralog

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
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

	containersList := []string{}
	for _, container := range a.getNodeContainers() {
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

	err = a.repository.InsertNode(node)
	if err != nil {
		log.Fatalln("couldn't insert node:", err)
	}
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}

func (a *App) validateURL(ans interface{}) error {
	val := reflect.ValueOf(ans).String()

	client, err := a.newClient(val)
	if err != nil {
		log.Fatalln("couldn't initiate client:", err)
	}

	err = a.pingServer()
	if err != nil {
		return errors.New("connection refused, please check your URL.")
	}

	a.centralogClient = client

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

	exists := a.repository.DoesNodeExist(val)
	if exists {
		return errors.New("a node with this name already exists, please choose a different name.")
	}

	return nil
}

func (a *App) getNodeContainers() []*pb.Container {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, _ := a.centralogClient.GetContainers(ctx, &pb.GetContainersRequest{})

	return response.Containers
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
