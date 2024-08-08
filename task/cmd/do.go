package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/apokryptein/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

// TODO: add ability to "do" multiple tasks
// e.g. -> task do 1 2 3
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark task as done",
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
		task, err := db.CompleteTask(taskId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error completing task in DB: %v", err)
			os.Exit(1)
		}

		fmt.Printf("You have completed the \"%s\" task.", task)
	},
}
