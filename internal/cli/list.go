package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/iv4n-t3a/wooordhunt-cli/internal/client"
)

type List struct {
	cli      *CLI
	tips     []client.Tips
	selected int
}

func newList(cli *CLI) List {
	return List{
		cli:      cli,
		tips:     nil,
		selected: -1,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.selected == -1 {
			m.selected = 0
		} else if msg.Type == tea.KeyCtrlJ && m.selected < len(m.tips)-1 {
			m.selected += 1
		} else if msg.Type == tea.KeyCtrlK && m.selected > 0 {
			m.selected -= 1
		}
		return m, nil
	}

	return m, nil
}

func (m List) View() (res string) {
	if m.tips == nil {
		return ""
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	// border := color.New(color.BgHiRed).SprintfFunc()
	maxlen := 0

	for i := range m.tips {
		maxlen = max(maxlen, len(m.tips[i].Word))
	}

	for i := range m.tips {
		word := m.tips[i].Word
		tips := m.tips[i].Tips
		spacesCount := maxlen - len(word) + 1
		spaces := strings.Repeat(" ", spacesCount)

		if i == m.selected {
			res += fmt.Sprintf("> %s: %s%s\n", cyan(word), spaces, tips)
		} else {
			res += fmt.Sprintf("%s: %s%s\n", cyan(word), spaces, tips)
		}
	}

	return
}

func (m *List) SetTips(tips []client.Tips) {
	if tips == nil {
		m.selected = -1
		m.tips = nil
		return
	}

	if m.selected >= len(tips) {
		m.selected = len(tips) - 1
	}
	m.tips = tips
}
