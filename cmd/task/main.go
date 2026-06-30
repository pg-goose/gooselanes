package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/pg-goose/gooselane/widgets/task"
)

func main() {
	if _, err := tea.NewProgram(task.New()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
