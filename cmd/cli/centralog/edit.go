package centralog

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
)

func (a *App) EditNodeCmd(nodeName string) {
	if !a.repository.DoesNodeExist(nodeName) {
		fmt.Println("a node with this name does not exist")
		return
	}

	selection := ""

	prompt := &survey.Select{
		Message: "Settings for " + nodeName,
		Options: []string{"Change node name", "Change target URL", "Change API key"},
	}

	survey.AskOne(prompt, &selection)

	switch selection {
	case "Change node name":
		a.editNodeNamePrompt(nodeName)
	case "Change target URL":
		a.editTargetURL(nodeName)
	case "Change API key":
		a.editAPIKey(nodeName)
	}
}

func (a *App) editAPIKey(nodeName string) {
	node, err := a.repository.GetNode(nodeName)
	if err != nil {
		log.Fatalln(err)
	}

	err = a.initClient(node.Location)
	if err != nil {
		log.Fatalln(err)
	}

	qs := []*survey.Question{
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "New API key:", Default: node.APIKey},
			Validate: a.validateKey,
		},
		{
			Name:   "confirm",
			Prompt: &survey.Confirm{Message: "Are you sure that you want to change this API key?"},
		},
	}

	var answers struct {
		Key     string
		Confirm bool
	}

	survey.Ask(qs, &answers)

	if !answers.Confirm {
		return
	}

	err = a.repository.UpdateAPIKey(nodeName, answers.Key)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("API key successfully updated")
	a.EditNodeCmd(nodeName)
}

func (a *App) editTargetURL(nodeName string) {
	node, err := a.repository.GetNode(nodeName)
	if err != nil {
		log.Fatalln(err)
	}

	qs := []*survey.Question{
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "New target URL:", Default: node.Location},
			Validate: a.validateURL,
		},
		{
			Name:   "confirm",
			Prompt: &survey.Confirm{Message: "Are you sure that you want to change this target URL?"},
		},
	}

	var answers struct {
		URL     string
		Confirm bool
	}

	survey.Ask(qs, &answers)

	if !answers.Confirm {
		return
	}

	err = a.repository.UpdateTargetURL(nodeName, answers.URL)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Target URL successfully updated")
	a.EditNodeCmd(nodeName)
}

func (a *App) editNodeNamePrompt(nodeName string) {
	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "New node name:", Default: nodeName},
			Validate: a.validateNodeName,
		},
		{
			Name:   "confirm",
			Prompt: &survey.Confirm{Message: "Are you sure that you want to rename this node?"},
		},
	}

	var answers struct {
		Name    string
		Confirm bool
	}

	survey.Ask(qs, &answers)

	if !answers.Confirm {
		return
	}

	err := a.repository.UpdateNodeName(nodeName, answers.Name)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Node successfully renamed")
	a.EditNodeCmd(answers.Name)
}
