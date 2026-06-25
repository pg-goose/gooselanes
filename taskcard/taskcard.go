package taskcard

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Surface int

const (
	SurfaceTitle Surface = iota
	SurfaceNotes
	SurfaceSteps
)

type stepItem struct {
	text string
	done bool
}

func (i stepItem) New(text string) stepItem {
	return stepItem{"text", false}
}

func (i stepItem) FilterValue() string {
	return i.text
}

func (i stepItem) Title() string {
	if i.done {
		return "[x] " + i.text
	}
	return "[ ] " + i.text
}

func (i stepItem) Description() string {
	return ""
}

type taskCardModel struct {
	surface Surface
	title   textinput.Model
	notes   textarea.Model
	steps   list.Model
}

func (card taskCardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			switch card.surface {
			case SurfaceTitle:
				card.surface = SurfaceSteps
			case SurfaceNotes:
				card.surface = SurfaceTitle
			case SurfaceSteps:
				if card.steps.Index() == 0 {
					card.surface = SurfaceNotes
					return card, card.notes.Focus()
				}
				card.steps.CursorUp()
			}
			return card, nil
		case "down":
			switch card.surface {
			case SurfaceTitle:
				card.surface = SurfaceNotes
			case SurfaceNotes:
				card.surface = SurfaceSteps
			case SurfaceSteps:
				last := len(card.steps.Items()) - 1
				if card.steps.Index() < last {
					card.steps.Select(card.steps.Index() + 1)
				}
			}
		case "enter":
			if card.surface == SurfaceSteps {
				index := card.steps.Index()
				if step, ok := card.steps.Items()[index].(stepItem); ok {
					step.done = !step.done
					return card, card.steps.SetItem(index, step)
				}
			}
		}
	}
	// Forward all other messages (typing, resize, etc.) to the active surface only
	var cmd tea.Cmd
	switch card.surface {
	case SurfaceTitle:
		card.title, cmd = card.title.Update(msg)
	case SurfaceNotes:
		card.notes, cmd = card.notes.Update(msg)
	case SurfaceSteps:
		card.steps, cmd = card.steps.Update(msg)
	}
	return card, cmd
}

func New() tea.Model {
	title := textinput.New()
	title.Placeholder = "Untitled"
	title.Focus()

	notes := textarea.New()
	notes.Placeholder = "Add notes..."

	steps := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 10)

	return taskCardModel{
		surface: SurfaceTitle,
		title:   title,
		notes:   notes,
		steps:   steps,
	}
}

// tea.Model requires these two alongside Update
func (card taskCardModel) Init() tea.Cmd {
	return textinput.Blink // start cursor blinking
}

func (card taskCardModel) View() tea.View {
	view := tea.NewView(fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		card.title.View(),
		card.notes.View(),
		card.steps.View(),
	))
	return view
}
