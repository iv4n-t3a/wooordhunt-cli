package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type input struct {
	cli       *CLI
	textInput textinput.Model
	err       error

	cachedWord string
	cachedView string
}

func newInput(cli *CLI) input {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Placeholder = ""

	return input{
		cli:        cli,
		textInput:  ti,
		err:        nil,
		cachedWord: "",
	}
}

func (m *input) Init() tea.Cmd {
	return textinput.Blink
}

func (m *input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *input) UpdateView() {

}

func (m *input) View() string {
	if len(m.textInput.Value()) == 0 {
		return fmt.Sprintf(
			"%s\n%s\n",
			m.textInput.View(),
			"(esc to quit)\n",
		)
	}

	if !strings.EqualFold(m.textInput.Value(), m.cachedWord) {
		tips, err := m.cli.client.GetTips(m.textInput.Value())

		if err != nil {
			m.cachedView = fmt.Sprintf(
				"%s\n%s\n",
				m.textInput.View(),
				err.Error(),
			)
		} else {
			m.cachedView = m.textInput.View() + "\n"
			maxlen := 0

			for i := range tips.Tips {
				maxlen = max(maxlen, len(tips.Tips[i].Word))
			}

			cyan := color.New(color.FgCyan).SprintFunc()

			for i := range tips.Tips {
				word := tips.Tips[i].Word
				tips := tips.Tips[i].Tips
				spacesCount := maxlen - len(word) + 1
				spaces := strings.Repeat(" ", spacesCount)
				m.cachedView += fmt.Sprintf("%s: %s%s\n", cyan(word), spaces, tips)
			}
		}
		m.cachedView += "(esc to quit)\n"
    m.cachedWord = m.textInput.Value()
	}

	return m.cachedView
}
