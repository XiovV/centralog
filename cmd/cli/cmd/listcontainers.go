/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

// listContainersCmd represents the containers command
var listContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "used to show containers for a specific node",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.ListContainersCmd(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	listCmd.AddCommand(listContainersCmd)
}
