/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

var (
	json bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "shows logs",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()
		app.ShowLogs()
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&json, "json", "j", false, "json output")

	rootCmd.AddCommand(logsCmd)
}
