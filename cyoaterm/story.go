package cyoaterm

import (
	"encoding/json"
	"fmt"
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
