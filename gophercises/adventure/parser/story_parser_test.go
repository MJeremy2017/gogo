package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStory(t *testing.T) {
	t.Run("Can parse json", func(t *testing.T) {
		fp := "../story.json"
		_, err := ParseStory(fp)
		assert.NoError(t, err, "parse story failed")
	})
}
