package cli

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

type input struct {
	cli       *CLI
	textInput textinput.Model
	spinner   spinner.Model
	err       error

	cachedView   string
	isUpdating   bool
	cancelUpdate func()
	mutex        sync.Mutex
}

func newInput(cli *CLI) input {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Placeholder = ""

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	return input{
		cli:       cli,
		textInput: ti,
		err:       nil,
		spinner:   sp,
	}
}

func (m *input) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m *input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			ctx, cancel := context.WithCancel(context.Background())
			if m.cancelUpdate != nil {
				m.cancelUpdate()
			}
			m.cancelUpdate = cancel
			m.isUpdating = true

			go func() {
				<-ctx.Done()
				m.UpdateView(m.textInput.Value())
			}()
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *input) UpdateView(word string) {
	tips, err := m.cli.client.GetTips(word)
	res := ""

	if err != nil {
		res = fmt.Sprintf(
			"%s\n%s\n",
			m.textInput.View(),
			err.Error(),
		)
	} else {
		maxlen := 0

		for i := range tips.Tips {
			maxlen = max(maxlen, len(tips.Tips[i].Word))
		}

		res += m.textInput.View() + "\n"
		cyan := color.New(color.FgCyan).SprintFunc()

		for i := range tips.Tips {
			word := tips.Tips[i].Word
			tips := tips.Tips[i].Tips
			spacesCount := maxlen - len(word) + 1
			spaces := strings.Repeat(" ", spacesCount)
			res += fmt.Sprintf("%s: %s%s\n", cyan(word), spaces, tips)
		}
	}
	res += "(esc to quit)\n"

  m.mutex.Lock()
  defer m.mutex.Unlock()

	m.cachedView = res
	m.isUpdating = false
}

func (m *input) View() string {
	if len(m.textInput.Value()) == 0 {
		return fmt.Sprintf(
			"%s\n%s\n",
			m.textInput.View(),
			"(esc to quit)\n",
		)
	}

	if m.isUpdating {
		return fmt.Sprintf(
			"%s\n%s\n%s\n",
			m.textInput.View(),
			m.spinner.View(),
			"(esc to quit)\n",
		)
	}

	return m.cachedView
}
