/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addNodeCmd represents the node command
var deleteNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "delete a node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0], "successfully deleted")
	},
}

func init() {
	deleteCmd.AddCommand(deleteNodeCmd)
}
