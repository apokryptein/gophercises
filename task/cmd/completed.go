package cmd

import (
	"fmt"
	"os"

	"github.com/apokryptein/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.ListCompleted()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error retrieving tasks: %v", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You haven't completed any todos today.")
			return
		}

		fmt.Println("You have finished the following tasks today:")
		for _, task := range tasks {
			fmt.Printf("- %s\n", task)
		}
	},
}
