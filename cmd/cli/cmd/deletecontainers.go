/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

// listContainersCmd represents the containers command
var deleteContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.DeleteContainersCmd(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	deleteCmd.AddCommand(deleteContainersCmd)
}
