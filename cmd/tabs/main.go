package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/pg-goose/gooselane/widgets/tabs"
)

func main() {
	titles := []string{
		"Tab 1", "Tab 2", "Tab 3", "Tab 4",
	}
	tabContent := []string{
		"Tab 1 Content", "Tab 2 Content", "Tab 3 Content", "Tab 4 Content",
	}
	m := tabs.New(titles, tabContent, true)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
