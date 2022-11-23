package parser

import (
	"encoding/json"
	"log"
	"os"
)

type Story struct {
	Intro struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	} `json:"intro"`
}

func ParseStory(fp string) {
	f, err := os.Open(fp)
	LogFatalErr(err)

	var b []byte // TODO: byte size can not be zero
	_, err = f.Read(b)
	LogFatalErr(err)

	var story Story
	err = json.Unmarshal(b, &story)
	LogFatalErr(err)
}

func LogFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
