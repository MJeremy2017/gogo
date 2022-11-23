package parser

import "testing"

func TestParseStory(t *testing.T) {
	t.Run("Can parse json", func(t *testing.T) {
		fp := "/Users/zhangyue/Workspace/github/gogo/gophercises/adventure/story.json"
		ParseStory(fp)

	})
}
