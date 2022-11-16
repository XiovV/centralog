/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

var (
	force bool
)

// addNodeCmd represents the node command
var deleteNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "delete a node",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.DeleteNodeCmd(args[0], force)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	deleteNodeCmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation prompt")

	deleteCmd.AddCommand(deleteNodeCmd)
}
