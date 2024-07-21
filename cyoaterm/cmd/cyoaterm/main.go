package main

import (
	"cyoaterm"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// brief server welcome text
	fmt.Println("Welcome to Choose Your Own Adventure")
	fmt.Println()

	// get flags
	filename := flag.String("f", "gopher.json", "name of CYOA JSON file")
	start := flag.String("s", "intro", "starting chapter of your story")
	tmpl := flag.String("t", "../../template.out", "desired story template")
	flag.Parse()

	// read in JSON file containing story data
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cyoa: error opening file %s: %v\n", *filename, err)
	}

	// parse JSON
	s, _ := cyoaterm.ParseJson(data)

	t := template.Must(template.ParseFiles(*tmpl))
	storyRepl(s, t, *start)
}

func storyRepl(s cyoaterm.Story, t *template.Template, chapter string) {
	if chapter == "home" {
		fmt.Println("we're going home")
		execTemplate(s, t, chapter)
		os.Exit(0)
	}

	execTemplate(s, t, chapter)

	for {
		fmt.Printf("Where are we going next: ")
		var next int
		fmt.Scanf("%d", &next)
		storyRepl(s, t, s[chapter].Options[next].Arc)
	}
}

func execTemplate(s cyoaterm.Story, t *template.Template, chapter string) {
	ClearScreen()
	err := t.Execute(os.Stdout, s[chapter])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cyoaterm: error executing template: %v\n", err)
	}
}

// Function taken from:
// https://dev.to/muhammadsaim/simplifying-your-terminal-experience-with-go-clearing-the-screen-1p7f
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
