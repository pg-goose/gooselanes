package task

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// notesModel is the task-card notes field. Focused: live textarea.
// Blurred: a static, dimmed box rendering the current text (no cursor).
type notesModel struct {
	area    textarea.Model
	focused bool
	width   int

	header lipgloss.Style // "Notes" section bar
	box    lipgloss.Style // unfocused static-box style
}

func newNotes() notesModel {
	ta := textarea.New()
	ta.Placeholder = "Add notes..."

	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	return notesModel{area: ta, header: header, box: box}
}

func (n *notesModel) Focus() tea.Cmd {
	n.focused = true
	return n.area.Focus()
}

func (n *notesModel) Blur() {
	n.focused = false
	n.area.Blur()
}

func (n notesModel) Focused() bool {
	return n.focused
}

// SetWidth expands the notes to fill the available space.
func (n *notesModel) SetWidth(w int) {
	n.width = w
	n.area.SetWidth(w)
}

func (n *notesModel) Value() string {
	return n.area.Value()
}

func (n *notesModel) SetValue(s string) {
	n.area.SetValue(s)
}

// Line and LineCount are surfaced for the parent's up/down surface switching.
func (n notesModel) Line() int {
	return n.area.Line()
}

func (n notesModel) LineCount() int {
	return n.area.LineCount()
}

func (n notesModel) Update(msg tea.Msg) (notesModel, tea.Cmd) {
	if !n.focused {
		return n, nil // static box is read-only; ignore input
	}
	var cmd tea.Cmd
	n.area, cmd = n.area.Update(msg)
	return n, cmd
}

func (n notesModel) View() string {
	header := n.header
	if n.width > 0 {
		header = header.Width(n.width)
	}
	bar := header.Render("Notes")

	if n.focused {
		return bar + "\n" + n.area.View()
	}

	text := n.area.Value()
	if text == "" {
		text = n.area.Placeholder
	}

	box := n.box
	if n.width > 0 {
		box = box.Width(n.width)
	}
	return bar + "\n" + box.Render(text)
}
