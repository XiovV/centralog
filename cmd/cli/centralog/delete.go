package centralog

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"strings"
)

func (a *App) DeleteContainersCmd(nodeName string) {
	if !a.repository.DoesNodeExist(nodeName) {
		fmt.Println("A node with this name does not exist")
		return
	}

	node, err := a.repository.GetNode(nodeName)
	if err != nil {
		log.Fatalln(err)
	}

	containers := strings.Split(node.Containers, ",")

	selectedContainers := []string{}
	prompt := &survey.MultiSelect{
		Message: "Which containers would you like to delete?",
		Options: containers,
	}
	survey.AskOne(prompt, &selectedContainers)

	var confirmationText string

	if len(selectedContainers) == 1 {
		confirmationText = "Are you sure that you want to delete this container?"
	} else {
		confirmationText = "Are you sure that you want to delete these containers?"
	}

	confirmation := a.confirmationPrompt(confirmationText)
	if !confirmation {
		return
	}

	remainingContainers := a.deleteContainers(containers, selectedContainers)

	err = a.repository.UpdateContainers(nodeName, strings.Join(remainingContainers, ","))
	if err != nil {
		log.Fatalln(err)
	}

	if len(selectedContainers) == 1 {
		fmt.Println("Container deleted successfully")
	} else {
		fmt.Println("Containers deleted successfully")
	}
}

func (a *App) deleteContainers(containers, toDelete []string) []string {
	for i, v := range containers {
		for _, v2 := range toDelete {
			if v == v2 {
				containers = append(containers[:i], containers[i+1:]...)
				break
			}
		}
	}

	return containers
}

func (a *App) DeleteNodeCmd(nodeName string, force bool) {
	if !a.repository.DoesNodeExist(nodeName) {
		fmt.Println("A node with this name does not exist")
		return
	}

	if force {
		err := a.repository.DeleteNode(nodeName)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Node deleted successfully")
		return
	}

	var confirmation bool
	prompt := &survey.Confirm{
		Message: "Are you sure that you want to delete " + nodeName + "?",
	}
	survey.AskOne(prompt, &confirmation)

	if !confirmation {
		return
	}

	a.repository.DeleteNode(nodeName)

	fmt.Println("Node deleted successfully")
}
