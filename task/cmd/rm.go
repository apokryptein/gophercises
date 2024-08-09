package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/apokryptein/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

// TODO: add ability to remove multiple tasks
// e.g. -> task rm 1 2 3
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove task from todo list",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error parsing task ID: %v", err)
			os.Exit(1)
		}

		tasks, err := db.ListTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error retrieving tasks: %v", err)
			os.Exit(1)
		}

		if id < 1 || id > len(tasks) {
			fmt.Fprint(os.Stderr, "task: invalid task ID")
			os.Exit(1)
		}

		taskId := tasks[id-1].Id
		task, err := db.RemoveTask(taskId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error removing task in DB: %v", err)
			os.Exit(1)
		}

		fmt.Printf("You have deleted the \"%s\" task.", task)
	},
}
