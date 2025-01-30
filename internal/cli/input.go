package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

func (m input) Init() tea.Cmd {
	return textinput.Blink
}

func (m input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m input) View() string {
	if strings.Compare(m.textInput.Value(), m.cachedWord) != 0 {
		tips, err := m.cli.client.GetTips(m.textInput.Value())

		if err != nil {
			m.cachedView = fmt.Sprintf(
				"%s \n%s \n%s \n",
				m.textInput.View(),
				err.Error(),
				"(esc to quit)",
			)
		} else {
			m.cachedView = fmt.Sprintf(
				"%s \n%s%s%s \n%s\n",
				m.textInput.View(),
				tips.Tips[0].Word, ": ", tips.Tips[0].Tips,
				"(esc to quit)",
			)
		}
	}

	return m.cachedView
}
