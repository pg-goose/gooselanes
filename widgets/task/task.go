package task

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Surface int

const (
	SurfaceTitle Surface = iota
	SurfaceNotes
	SurfaceSteps
)

type Model struct {
	surface   Surface
	title     model
	notes     notesModel
	checklist checklistModel
}

// focus blurs every surface, focuses the active one, returns its blink cmd.
func (card *Model) focus(s Surface) tea.Cmd {
	card.surface = s
	card.title.Blur()
	card.notes.Blur()
	card.checklist.Blur()
	switch s {
	case SurfaceTitle:
		return card.title.Focus()
	case SurfaceNotes:
		return card.notes.Focus()
	case SurfaceSteps:
		return card.checklist.Focus()
	}
	return nil
}

func (card Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		card.title.SetWidth(msg.Width)
		card.notes.SetWidth(msg.Width)
		card.checklist.SetWidth(msg.Width)
		// then drops to the forward block, which relays the resize to the active surface
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			switch card.surface {
			case SurfaceTitle:
				return card, card.focus(SurfaceSteps)
			case SurfaceNotes:
				if card.notes.Line() == 0 {
					return card, card.focus(SurfaceTitle)
				}
			case SurfaceSteps:
				// only leave the checklist when sitting on its first item;
				// otherwise the checklist handles up internally (forwarded below).
				if card.checklist.AtTop() {
					return card, card.focus(SurfaceNotes)
				}
			}
		case "down":
			switch card.surface {
			case SurfaceTitle:
				return card, card.focus(SurfaceNotes)
			case SurfaceNotes:
				if card.notes.Line() >= card.notes.LineCount()-1 {
					return card, card.focus(SurfaceSteps)
				}
				// SurfaceSteps: down never leaves; forwarded to the checklist below.
			}
		}
	}
	// forward all other messages (typing, resize, etc.) to the active surface only
	var cmd tea.Cmd
	switch card.surface {
	case SurfaceTitle:
		card.title, cmd = card.title.Update(msg)
	case SurfaceNotes:
		card.notes, cmd = card.notes.Update(msg)
	case SurfaceSteps:
		card.checklist, cmd = card.checklist.Update(msg)
	}
	return card, cmd
}

func New() Model {
	title := newTitle()
	title.Focus()

	notes := newNotes()

	checklist := newChecklist()

	return Model{
		surface:   SurfaceTitle,
		title:     title,
		notes:     notes,
		checklist: checklist,
	}
}

func (card Model) Init() tea.Cmd {
	return card.title.Blink() // start cursor blinking
}

func (card Model) View() tea.View {
	view := tea.NewView(fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		card.title.View(),
		card.notes.View(),
		card.checklist.View(),
	))
	return view
}
