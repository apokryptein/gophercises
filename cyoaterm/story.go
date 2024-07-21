package cyoaterm

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime"
)

type Story map[string]Arc

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func ParseJson(data []byte) (Story, error) {
	var storyData Story
	if err := json.Unmarshal(data, &storyData); err != nil {
		return nil, err
	}
	return storyData, nil
}

func (s Story) PrintArcTitles() {
	for k := range s {
		fmt.Printf("%s: %s\n", k, s[k].Title)
	}
}

func StoryRepl(s Story, t *template.Template, chapter string) {
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
		StoryRepl(s, t, s[chapter].Options[next].Arc)
	}
}

func execTemplate(s Story, t *template.Template, chapter string) {
	clearScreen()
	err := t.Execute(os.Stdout, s[chapter])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cyoaterm: error executing template: %v\n", err)
	}
}

// Function taken from:
// https://dev.to/muhammadsaim/simplifying-your-terminal-experience-with-go-clearing-the-screen-1p7f
func clearScreen() {
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
