package client_test

import (
	"os"
	"testing"

	"github.com/iv4n-t3a/wooordhunt-cli/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestEnHtmlToWordInfo(t *testing.T) {
	htmlText, err := os.ReadFile("../../data/apple.html")
	assert.NoError(t, err, "Can't open test data file")

	wordInfo := client.HtmlToWordInfo(string(htmlText), "apple")

	assert.Equal(t, "apple", wordInfo.Word)
	assert.Equal(t, "яблоко, яблоня, чепуха, лесть, яблочный", *wordInfo.Meaning)
	assert.Equal(t, "существительное ↓", *wordInfo.WordType)
	assert.Equal(t, 12, len(wordInfo.Phrases))
	assert.Equal(t, "the lone ripe apple in the entire bag\u2002—\u2002одно-единственное спелое яблоко в целом мешке", wordInfo.Phrases[0])
	assert.Equal(t, 0, len(wordInfo.Examples))
	assert.Equal(t, []string{
		"applied \u2002—\u2002прикладной, приложенный",
		"apply \u2002—\u2002применять, относиться, использовать, обращаться, прикладывать, касаться, прилагать",
	}, wordInfo.SimilarWords)
}
