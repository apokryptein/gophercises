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

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove task from todo list",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		for _, v := range args {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "task: error parsing task ID: %v", err)
				os.Exit(1)
			}
			ids = append(ids, id)
		}

		tasks, err := db.ListTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error retrieving tasks: %v", err)
			os.Exit(1)
		}

		for _, id := range ids {
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

			fmt.Printf("You have deleted the \"%s\" task.\n", task)
		}
	},
}
