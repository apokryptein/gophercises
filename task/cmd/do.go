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

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark task as done",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error parsing task ID: %v", err)
			os.Exit(1)
		}

		task, err := db.CompleteTask(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task: error completing task in DB: %v", err)
			os.Exit(1)
		}

		fmt.Printf("You have completed the \"%s\" task.", task)
	},
}
