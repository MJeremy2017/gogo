package parser

import (
	"encoding/json"
	"os"
)

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func ParseStory(fp string) (Story, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return Story{}, err
	}

	var story Story
	err = json.Unmarshal(b, &story)
	if err != nil {
		return Story{}, err
	}
	return story, nil
}
