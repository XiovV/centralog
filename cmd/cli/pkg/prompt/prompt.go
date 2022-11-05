package prompt

import (
	"errors"
	"github.com/manifoldco/promptui"
)

func GetURL(label string) string {
	//TODO: add url validation here
	prompt := promptui.Prompt{
		Label: label,
	}

	result, _ := prompt.Run()
	return result
}

func GetAPIKey(label string) (string, error) {
	validate := func(input string) error {
		if len(input) != 5 {
			return errors.New("api key must be exactly 5 characters long")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetString(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, _ := prompt.Run()
	return result
}
