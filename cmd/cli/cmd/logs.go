/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/XiovV/centralog-agent/cmd/cli/centralog"

	"github.com/spf13/cobra"
)

var (
	containers []string
	showAll    bool
	follow     bool
	first      int32
	last       int32
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "shows logs",
	Run: func(cmd *cobra.Command, args []string) {
		app := centralog.NewApp()

		nodeName := args[0]
		flags := centralog.ShowLogsFlags{
			Containers: containers,
			ShowAll:    showAll,
			First:      first,
			Last:       last,
			Follow:     true,
		}

		if first > 0 || last > 0 {
			flags.Follow = false
		}

		if follow {
			flags.Follow = true
		}

		app.ShowLogs(nodeName, flags)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	logsCmd.Flags().StringSliceVarP(&containers, "containers", "c", []string{}, "exampleContainer1,exampleContainer2,exampleContainer3")
	logsCmd.Flags().BoolVarP(&showAll, "all", "a", false, "shows all logs")
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "receive logs in real time")
	logsCmd.Flags().Int32Var(&first, "first", 0, "shows first n amount of logs")
	logsCmd.Flags().Int32Var(&last, "last", 0, "shows last n amount of logs")

	rootCmd.AddCommand(logsCmd)
}
