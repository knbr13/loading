package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knbr13/glow-net/internal/enrichment"
	"github.com/knbr13/glow-net/internal/state"
	"github.com/knbr13/glow-net/internal/ui"
)

func main() {
	enricher, err := enrichment.NewEnricher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing enricher: %v\n", err)
	}
	defer enricher.Close()

	s := state.NewAppState(enricher)
	m := ui.NewModel(s)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
