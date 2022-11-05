/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"
	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "gives you a list of all of your nodes",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.App{}
		app.ListNodesCmd()
	},
}

func init() {
	listCmd.AddCommand(nodesCmd)
}
