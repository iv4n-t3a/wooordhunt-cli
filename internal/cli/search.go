package cli

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

type Search struct {
	cli       *CLI
	textInput textinput.Model
	spinner   spinner.Model
	tips      List
	err       error

	isUpdating   bool
	lastUpdateId int
	mutex        sync.Mutex
}

func newSearch(cli *CLI) Search {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Placeholder = ""

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	return Search{
		cli:       cli,
		textInput: ti,
		spinner:   sp,
		tips:      newList(cli),
		err:       nil,
	}
}

func (m *Search) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m *Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.Type == tea.KeyEnter && m.tips.selected != -1 {
			return m.OpenWordInfo()
		}
		if t := msg.Type; t == tea.KeyCtrlJ || t == tea.KeyCtrlK {
			m.tips, cmd = m.tips.Update(msg)
		} else if t := msg.Type; t == tea.KeyCtrlC || t == tea.KeyEscape {
			return m, tea.Quit
		} else {
			m.textInput, cmd = m.textInput.Update(msg)
			m.lastUpdateId++
			go m.UpdateTips(m.textInput.Value(), m.lastUpdateId)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case error:
		m.err = msg
		return m, nil
	}
	return m, cmd
}

func (m *Search) View() string {
	if len(m.textInput.Value()) == 0 {
		return fmt.Sprintf(
			"%s\n",
			m.textInput.View(),
		)
	}

	if m.isUpdating {
		return fmt.Sprintf(
			"%s\n%s\n",
			m.textInput.View(),
			m.spinner.View(),
		)
	}

	return fmt.Sprintf(
		"%s\n%s\n",
		m.textInput.View(),
		m.tips.View(),
	)
}

func (m *Search) UpdateTips(word string, updateId int) {
	m.isUpdating = true
	tips, err := m.cli.client.GetTips(word)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if updateId != m.lastUpdateId {
		return
	}

	if err != nil {
		m.tips.SetTips(nil)
	} else {
		m.tips.SetTips(tips.Tips)
	}

	m.isUpdating = false
}

func (m *Search) OpenWordInfo() (tea.Model, tea.Cmd) {
	word := m.tips.tips[m.tips.selected].Word
	wordInfo, err := m.cli.client.GetWord(m.textInput.Value())

	if err != nil {
		return newPager("Error", err.Error(), m.cli, m), nil
	}

	bold := color.New(color.Bold).SprintFunc()
	italic := color.New(color.Italic).SprintFunc()
	grey := color.RGB(100, 100, 100).SprintFunc()

	text := fmt.Sprintf(
		"\n\nЧасть речи: %s\n\nЗначение: %s\n\n",
		grey(italic(wordInfo.WordType)),
		bold(wordInfo.Meaning),
	)

	text += "Словосочетания:\n"
	text = AddListOfPhrases(wordInfo.Phrases, text)

	text += "Похожие слова:\n"
	text = AddListOfPhrases(wordInfo.SimilarWords, text)

	text += "Формы слова:\n"
	text = AddListOfPhrases(wordInfo.WordForms, text)

	return newPager(text, word, m.cli, m), nil
}

func AddListOfPhrases(list []string, text string) string {
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
