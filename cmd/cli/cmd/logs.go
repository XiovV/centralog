/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

var (
	json       bool
	containers []string
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "shows logs",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.ShowLogs(containers)
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&json, "json", "j", false, "json output")
	logsCmd.Flags().StringSliceVarP(&containers, "containers", "c", []string{}, "exampleContainer1,exampleContainer2,exampleContainer3")

	rootCmd.AddCommand(logsCmd)
}
