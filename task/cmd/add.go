package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/apokryptein/gophercises/task/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Task to Queue",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.AddTask(task)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error adding task: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Added \"%s\" to your todo list.", task)
	},
}
