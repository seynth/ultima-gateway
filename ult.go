package main

import (
	"fmt"
	"os"
	"ultima/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.UltimaInit(), tea.WithAltScreen())
	ultima, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if _, yes := ultima.(ui.Ultima); yes {
		fmt.Println("Thank you for using Ultima Gateway :-)")
	}
}
