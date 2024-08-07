package cmd

import (
	"fmt"
	"os"

	"github.com/apokryptein/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.ListTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error retrieving tasks: %v", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You do not have pending todos")
		}

		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.Id, task.Name)
		}
	},
}
