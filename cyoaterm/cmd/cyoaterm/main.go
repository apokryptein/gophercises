package main

import (
	"cyoaterm"
	"flag"
	"fmt"
	"html/template"
	"os"
)

func main() {
	// brief server welcome text
	fmt.Println("Welcome to Choose Your Own Adventure")
	fmt.Println()

	// get flags
	filename := flag.String("f", "../../configs/stories/gopher.json", "name of CYOA JSON file")
	start := flag.String("s", "intro", "starting chapter of your story")
	tmpl := flag.String("t", "../../configs/templates/default.tmpl", "desired story template")
	flag.Parse()

	// read in JSON file containing story data
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cyoa: error opening file %s: %v\n", *filename, err)
	}

	// parse JSON
	s, _ := cyoaterm.ParseJson(data)

	t := template.Must(template.ParseFiles(*tmpl))
	cyoaterm.StoryRepl(s, t, *start)
}
