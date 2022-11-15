package centralog

import (
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
