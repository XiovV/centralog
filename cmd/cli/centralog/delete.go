package centralog

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
)

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
