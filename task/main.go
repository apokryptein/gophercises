package main

import (
	"github.com/apokryptein/gophercises/task/cmd"
	"github.com/apokryptein/gophercises/task/db"
)

func main() {
	db.Init()
	cmd.Execute()
}
