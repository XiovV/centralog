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
		Options: []string{"Change node name", "Change target URL", "Change API key", "Change containers"},
	}

	survey.AskOne(prompt, &selection)

	switch selection {
	case "Change node name":
		a.editNodeNamePrompt(nodeName)
	}
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
}
