/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-cli/centralog"

	"github.com/spf13/cobra"
)

// containersCmd represents the containers command
var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "used to show containers for a specific node",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.App{}
		app.ListContainersCmd(args[0])
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	listCmd.AddCommand(containersCmd)
}
