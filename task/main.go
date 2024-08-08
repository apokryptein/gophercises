package main

import (
	"github.com/apokryptein/gophercises/task/cmd"
	"github.com/apokryptein/gophercises/task/db"
)

// TODO: Add command line argument for desired DB file location
// default to ~/.config/task directory
func main() {
	db.Init()
	cmd.Execute()
}
