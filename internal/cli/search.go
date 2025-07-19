package cli

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			return m.openWordInfo()
		}
		if t := msg.Type; t == tea.KeyDown || t == tea.KeyUp {
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

func (m *Search) openWordInfo() (tea.Model, tea.Cmd) {
	word := m.tips.tips[m.tips.selected].Word
	wordInfo, err := m.cli.client.GetWord(word)

	if err != nil {
    return getErrorPager(err, m.cli, m), nil
	}

  return getWordPager(*wordInfo, m.cli, m), nil
}
