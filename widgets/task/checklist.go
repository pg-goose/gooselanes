package task

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type stepItem struct {
	text string
	done bool
}

func (i stepItem) FilterValue() string {
	return i.text
}

func (i stepItem) Title() string {
	if i.done {
		return "- [x] " + i.text
	}
	return "- [ ] " + i.text
}

func (i stepItem) Description() string {
	return ""
}

// checklistModel is the task-card checklist. Focused: navigate items with
// up/down; the bottom row is a text input that appends new items. Blurred:
// a static slate panel. The whole panel sits on a distinct background color.
type checklistModel struct {
	list    list.Model
	input   textinput.Model
	focused bool
	onInput bool // sub-focus: true = add-input, false = list items
	width   int

	header lipgloss.Style // "Checklist" section bar
	panel  lipgloss.Style // slate background wrapping list + input
}

func newChecklist() checklistModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 10)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "add item..."
	ti.Prompt = "- [ ] "

	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#2E3650")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	panel := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E2230")).
		Padding(0, 1)

	return checklistModel{list: l, input: ti, header: header, panel: panel}
}

func (c *checklistModel) Focus() tea.Cmd {
	c.focused = true
	c.onInput = false // always land on the list first
	c.input.Blur()
	return nil
}

func (c *checklistModel) Blur() {
	c.focused = false
	c.onInput = false
	c.input.Blur()
}

func (c checklistModel) Focused() bool {
	return c.focused
}

// AtTop reports whether up should bubble out of the checklist (to the surface
// above). True only when sitting on the first list item, never while in the
// add-input — there, up moves back into the list instead.
func (c checklistModel) AtTop() bool {
	return !c.onInput && c.list.Index() == 0
}

func (c *checklistModel) SetWidth(w int) {
	c.width = w
	// panel padding (0,1) eats 2 cols; keep inner widgets inside it.
	inner := w - 2
	if inner < 0 {
		inner = 0
	}
	c.list.SetWidth(inner)
	c.input.SetWidth(inner - lipgloss.Width(c.input.Prompt))
}

// toggle flips the checkbox of the focused list item.
func (c *checklistModel) toggle() tea.Cmd {
	i := c.list.Index()
	items := c.list.Items()
	if i < 0 || i >= len(items) {
		return nil
	}
	if step, ok := items[i].(stepItem); ok {
		step.done = !step.done
		return c.list.SetItem(i, step)
	}
	return nil
}

func (c checklistModel) Update(msg tea.Msg) (checklistModel, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "up":
			if c.onInput {
				// leave input, drop onto the last list item
				c.onInput = false
				c.input.Blur()
				if last := len(c.list.Items()) - 1; last >= 0 {
					c.list.Select(last)
				}
				return c, nil
			}
			// at the first item the parent handles the surface switch (AtTop);
			// otherwise move up one.
			if c.list.Index() > 0 {
				c.list.CursorUp()
			}
			return c, nil
		case "down":
			if c.onInput {
				return c, nil // already the bottom row
			}
			if last := len(c.list.Items()) - 1; c.list.Index() >= last {
				c.onInput = true
				return c, c.input.Focus()
			}
			c.list.CursorDown()
			return c, nil
		case "enter":
			if c.onInput {
				text := strings.TrimSpace(c.input.Value())
				if text != "" {
					cmd := c.list.InsertItem(len(c.list.Items()), stepItem{text: text})
					c.input.SetValue("") // clear to add another, stay on input
					return c, cmd
				}
				return c, nil
			}
			return c, c.toggle()
		case " ":
			if !c.onInput {
				return c, c.toggle()
			}
			// in the input, space types normally — fall through
		}
	}

	// forward everything else to whichever sub-widget is active
	var cmd tea.Cmd
	if c.onInput {
		c.input, cmd = c.input.Update(msg)
	} else {
		c.list, cmd = c.list.Update(msg)
	}
	return c, cmd
}

func (c checklistModel) View() string {
	header := c.header
	if c.width > 0 {
		header = header.Width(c.width)
	}
	bar := header.Render("Checklist")

	body := c.list.View() + "\n" + c.input.View()

	panel := c.panel
	if c.width > 0 {
		panel = panel.Width(c.width)
	}
	return bar + "\n" + panel.Render(body)
}
