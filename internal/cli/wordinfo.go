package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/iv4n-t3a/wooordhunt-cli/internal/client"
)

func getWordPager(info client.WordInfo, cli *CLI, model tea.Model) tea.Model {
	text := getWordDescription(info)
	return newPager(text, info.Word, cli, model)
}

func getErrorPager(err error, cli *CLI, model tea.Model) tea.Model {
	return newPager("Error", err.Error(), cli, model)
}

func getWordDescription(info client.WordInfo) string {
	bold := color.New(color.Bold).SprintFunc()
	italic := color.New(color.Italic).SprintFunc()
	grey := color.RGB(100, 100, 100).SprintFunc()

	text := "\n\n"

	if info.Meaning != nil {
		text += fmt.Sprintf(
      "Значение: %s\n\n",
			bold(*info.Meaning),
		)
	}

  if info.WordType != nil {
		text += fmt.Sprintf(
			"Часть речи: %s\n\n",
			grey(italic(*info.WordType)),
		)
	}

	if info.Phrases != nil {
		text += "Словосочетания:\n"
		text = addListOfPhrases(info.Phrases, text)
	}

	if info.SimilarWords != nil {
		text += "Похожие слова:\n"
		text = addListOfPhrases(info.SimilarWords, text)
	}

	if info.WordForms != nil {
		text += "Формы слова:\n"
		text = addListOfPhrases(info.WordForms, text)
	}

	return text
}

func addListOfPhrases(list []string, text string) string {
	italic := color.New(color.Italic).SprintFunc()
	grey := color.RGB(100, 100, 100).SprintFunc()

	for i := range list {
		phrase := list[i]

		separator := "—"
		if strings.Contains(phrase, separator) {
			parts := strings.Split(phrase, separator)

			text += fmt.Sprintf(
				"    %s%s%s\n",
				italic(parts[0]),
				separator,
				grey(italic(parts[1])),
			)
		} else {
			text += fmt.Sprintf(
				"    %s\n",
				italic(phrase),
			)
		}
	}
	return text + "\n"
}
