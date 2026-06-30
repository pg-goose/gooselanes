package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pg-goose/gooselane/taskcard"
	"github.com/pg-goose/gooselane/util"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type placeholderModel struct {
	area textarea.Model
}

func newPlaceholder(placeholder string) placeholderModel {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.Focus()
	return placeholderModel{area: ta}
}

func (p placeholderModel) Init() tea.Cmd  { return textarea.Blink }
func (p placeholderModel) View() tea.View { return tea.NewView(p.area.View()) }
func (p placeholderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	p.area, cmd = p.area.Update(msg)
	return p, cmd
}

type model struct {
	tabs      []tab
	activeTab int
}

type tab struct {
	title   string
	content tea.Model
}

func (m model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.tabs))
	for i, t := range m.tabs {
		cmds[i] = t.content.Init()
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+left":
			m.activeTab = util.PrevIndex(m.tabs, m.activeTab)
			return m, nil
		case "ctrl+right":
			m.activeTab = util.NextIndex(m.tabs, m.activeTab)
			return m, nil
		}
	}
	// Forward all other messages to active tab's content
	var cmd tea.Cmd
	m.tabs[m.activeTab].content, cmd = m.tabs[m.activeTab].content.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	doc := strings.Builder{}
	s := lipgloss.NewStyle()
	for i, tab := range m.tabs {
		if i > 0 {
			doc.WriteString("  |  ")
		}
		if i == m.activeTab {
			fmt.Fprintf(&doc, "<%s>", tab.title)
		} else {
			doc.WriteString(tab.title)
		}
	}
	doc.WriteString("\n")
	doc.WriteString(s.Render(m.tabs[m.activeTab].content.View().Content))
	v := tea.NewView(doc.String())
	v.AltScreen = true
	return v
}

func main() {
	m := model{
		tabs: []tab{
			{title: "Tasks", content: taskcard.New()},
			{title: "Checklists", content: newPlaceholder("Write notes...")},
			{title: "Notes", content: newPlaceholder("Board coming soon...")},
		},
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
