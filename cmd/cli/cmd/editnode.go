/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Opens a prompt for editing a node's settings",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.EditNodeCmd(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	editCmd.AddCommand(nodeCmd)
}
