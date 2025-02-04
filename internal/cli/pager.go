package cli

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Pager struct {
	cli      *CLI
	title    string
	viewport viewport.Model

	parrent tea.Model
}

func newPager(content string, title string, cli *CLI, parrent tea.Model) (m Pager) {
	width, height, err := term.GetSize(0)

	if err != nil {
		log.Fatal("Failed to get terminal size")
	}

	m.title = title
	m.cli = cli
	m.parrent = parrent
	m.viewport = viewport.New(width, height-m.verticalMarginHeight())
	m.viewport.YPosition = m.headerHeight()
	m.viewport.SetContent(content)
	return
}

func (m Pager) Init() tea.Cmd {
	return nil
}

func (m Pager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "backspace" || k == "esc" {
			if m.parrent != nil {
				return m.parrent, nil
			}
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)

			if m.parrent == nil {
				return m, tea.Quit
			}

			return m.parrent, tea.Batch(cmds...)
		}
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.verticalMarginHeight()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Pager) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m Pager) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Pager) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Pager) headerHeight() int {
	return lipgloss.Height(m.headerView())
}

func (m Pager) footerHeight() int {
	return lipgloss.Height(m.footerView())
}

func (m Pager) verticalMarginHeight() int {
	return m.headerHeight() + m.footerHeight()
}
