package tabs

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	doc         lipgloss.Style
	highlight   lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	window      lipgloss.Style
}

func newStyles(bgIsDark bool) *styles {
	lightDark := lipgloss.LightDark(bgIsDark)

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	highlightColor := lightDark(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))

	s := new(styles)
	s.doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	return s
}

type Tabs struct {
	Tabs       []string
	TabContent []string
	styles     *styles
	activeTab  int
}

func (m Tabs) Init() tea.Cmd {
	return nil
}

func (m Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+right":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "ctrl+left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}
	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m Tabs) View() tea.View {
	if m.styles == nil {
		return tea.NewView("")
	}
	doc := strings.Builder{}
	s := m.styles

	var renderedTabs []string

	for i, tab := range m.Tabs {
		var style lipgloss.Style
		isFirst := (i == 0)
		isLast := (i == len(m.Tabs)-1)
		isActive := (i == m.activeTab)

		style = s.inactiveTab
		if isActive {
			style = s.activeTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.window.Width((lipgloss.Width(row))).Render(m.TabContent[m.activeTab]))
	return tea.NewView(s.doc.Render(doc.String()))
}

func (m Tabs) SetContent(tab int, content string) {
	m.TabContent[tab] = content
}

func (m Tabs) SetTitle(tab int, title string) {
	m.Tabs[tab] = title
}

func New(tabs []string, content []string, isBgDark bool) Tabs {
	return Tabs{
		Tabs:       tabs,
		TabContent: content,
		styles:     newStyles(isBgDark),
	}
}
