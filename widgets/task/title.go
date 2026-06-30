package task

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// model is the task-card title field. Focused: editable textinput.
// Blurred: a solid section bar showing the current text.
type model struct {
	input   textinput.Model
	focused bool
	width   int

	bar lipgloss.Style // unfocused section-bar style
}

func newTitle() model {
	ti := textinput.New()
	ti.Placeholder = "Untitled"

	bar := lipgloss.NewStyle().
		Background(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	return model{input: ti, bar: bar}
}

func (t *model) Focus() tea.Cmd {
	t.focused = true
	return t.input.Focus()
}

func (t *model) Blur() {
	t.focused = false
	t.input.Blur()
}

func (t model) Focused() bool {
	return t.focused
}

// cursor-blink cmd surfaced so the parent's Init can start it.
func (t model) Blink() tea.Cmd {
	return textinput.Blink
}

// SetWidth expands the title to fill the available space.
func (t *model) SetWidth(w int) {
	t.width = w
	// textinput counts the prompt
	inner := w - lipgloss.Width(t.input.Prompt)
	if inner < 0 {
		inner = 0
	}
	t.input.SetWidth(inner)
}

func (t *model) Value() string {
	return t.input.Value()
}

func (t *model) SetValue(s string) {
	t.input.SetValue(s)
}

func (t model) Update(msg tea.Msg) (model, tea.Cmd) {
	if !t.focused {
		return t, nil // bar is read-only; ignore input
	}
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	return t, cmd
}

func (t model) View() string {
	if t.focused {
		return t.input.View()
	}
	text := t.input.Value()
	if text == "" {
		text = t.input.Placeholder
	}
	style := t.bar
	if t.width > 0 {
		style = style.Width(t.width)
	}
	return style.Render(text)
}
