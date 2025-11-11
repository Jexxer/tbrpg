package storage

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
)

// RenderHelpModal renders the storage help modal
func RenderHelpModal() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center)

	help := `
Storage Keybinds

Navigation:
  Tab         - Cycle focus
  ↑/↓ or j/k  - Navigate lists/table
  ←/→ or h/l  - Switch between categories/table
  Enter       - Select item

Search:
  /           - Activate search
  ESC         - Exit search

Search Syntax:
  text        - Match item name
  tag:weapon  - Filter by tag
  *ore        - Wildcard (ends with "ore")
  qty:>50     - Quantity greater than 50

Actions:
  S           - Save current search
  O           - Load saved search
  ?           - Show this help

Press ESC or ? to close
`

	return titleStyle.Render("Storage Help") + "\n" + help
}

// RenderSaveSearchModal renders the save search modal
func RenderSaveSearchModal(currentQuery, modalInputView string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	title := titleStyle.Render("Save Search")

	currentQueryText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Query: " + currentQuery)

	prompt := "Name: " + modalInputView

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("\nEnter = Save  |  ESC = Cancel")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		currentQueryText,
		"",
		prompt,
		instructions,
	)
}

// RenderLoadSearchModal renders the load search modal
func RenderLoadSearchModal(savedSearches []game.SavedSearch) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	title := titleStyle.Render("Saved Searches")

	var searches string
	if len(savedSearches) == 0 {
		searches = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("\nNo saved searches yet\n\nPress 's' while searching to save a search")
	} else {
		for i, search := range savedSearches {
			searches += fmt.Sprintf("\n%d. %s\n   %s\n", i+1, search.Name, search.Query)
		}
		searches += lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("\nPress number to load  |  ESC to close")
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, searches)
}
