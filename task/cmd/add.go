package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Task to Queue",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is a fake \"add\" command")
	},
}
