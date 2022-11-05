package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
)

type Node struct {
	URL  string
	Key  string
	Name string
}

func AddNode() (Node, []string) {
	qs := []*survey.Question{
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "Enter your node's URL:"},
			Validate: survey.Required,
		},
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Enter your node's API key:"},
			Validate: survey.Required,
		},
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Enter your node's custom name:"},
			Validate: survey.Required,
		},
	}

	var answers Node
	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatalln(err)
	}

	containers := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select containers:",
		Options: []string{"logserver1", "logserver2"},
	}

	survey.AskOne(prompt, &containers)

	return answers, containers
}

//func GetURL(label string) string {
//	//TODO: add url validation here
//	prompt := promptui.Prompt{
//		Label: label,
//	}
//
//	result, _ := prompt.Run()
//	return result
//}
//
//func GetAPIKey(label string) (string, error) {
//	validate := func(input string) error {
//		if len(input) != 5 {
//			return errors.New("api key must be exactly 5 characters long")
//		}
//
//		return nil
//	}
//
//	prompt := promptui.Prompt{
//		Label:    label,
//		Validate: validate,
//	}
//
//	result, err := prompt.Run()
//	if err != nil {
//		return "", err
//	}
//
//	return result, nil
//}
//
//func GetString(label string) string {
//	prompt := promptui.Prompt{
//		Label: label,
//	}
//
//	result, _ := prompt.Run()
//	return result
//}
//
//func ShowSelect() {
//	prompt := promptui.Select{
//		Label: "Select Containers",
//		Items: []string{"container1", "container2", "container3"},
//	}
//
//	_, result, err := prompt.Run()
//	if err != nil {
//		return err
//	}
//
//}
