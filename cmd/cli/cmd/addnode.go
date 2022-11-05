/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-cli/centralog"
	"github.com/spf13/cobra"
)

var (
	url, apiKey, name string
)

// addNodeCmd represents the node command
var addNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "used to add a node",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.App{}

		if url == "" && apiKey == "" && name == "" {
			app.AddNodeWithPrompt()
			return
		}

		app.AddNodeWithFlags(url, apiKey, name)
	},
}

func init() {
	addNodeCmd.Flags().StringVarP(&url, "url", "u", "", "url where your node is")
	addNodeCmd.Flags().StringVarP(&apiKey, "key", "k", "", "your api key")
	addNodeCmd.Flags().StringVarP(&name, "name", "n", "", "name for your node")

	addCmd.AddCommand(addNodeCmd)
}
