package main

import (
	"fmt"
	"os"

	"github.com/apokryptein/gophercises/task/cmd"
	"github.com/apokryptein/gophercises/task/db"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "task: error checking for user's home directory: %v", err)
		os.Exit(1)
	}

	dbPath := homeDir + "/.config/task/"
	if err := setupPath(dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "task: error creating DB directory: %v", err)
		os.Exit(1)
	}

	db.Init(dbPath)
	cmd.Execute()
}

// Creates directory ~/config/task if it doesn't exist
func setupPath(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}
