// Package views
package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
)

// StorageViewParams contains all the data needed to render the storage view
type StorageViewParams struct {
	Width               int
	Height              int
	FocusedColor        string
	StorageSearchActive bool
	StorageFocus        int // 0 = category, 1 = table

	// Component views (already rendered)
	SearchInputView  string
	CategoryListView string
	TableView        string

	// Table configuration
	AvailableWidth  int
	AvailableHeight int
}

const (
	StorageFocusCategory = 0
	StorageFocusTable    = 1
)

// RenderStorageView renders the storage view with category list and items table
func RenderStorageView(params StorageViewParams) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	topBarHeight := 3
	commandLineHeight := 3
	activityLogHeight := 7
	storageSearchBarHeight := 3

	// Calculate available dimensions
	availableWidth := params.Width - 15 - 25 - 10
	availableHeight := params.Height - topBarHeight - commandLineHeight - activityLogHeight - storageSearchBarHeight

	// Search bar
	searchBarStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Padding(0, 1)

	var searchBar string
	if params.StorageSearchActive {
		searchBar = searchBarStyle.Render("Search: " + params.SearchInputView)
	} else {
		searchBar = searchBarStyle.Render("Search: " + params.SearchInputView + " (press / to search)")
	}

	// Layout: Category list on left, table on right
	categoryWidth := 20
	tableWidth := availableWidth - categoryWidth + 1

	// Category panel with focus indicator
	categoryStyle := lipgloss.NewStyle().
		Width(categoryWidth).
		Height(len(game.GetCategories()))

	// Add border to show focus
	if params.StorageFocus == StorageFocusCategory {
		categoryStyle = categoryStyle.
			Border(lipgloss.RoundedBorder()).
			BorderRight(false).
			BorderForeground(lipgloss.Color(params.FocusedColor))
	} else {
		categoryStyle = categoryStyle.
			Border(lipgloss.RoundedBorder()).
			BorderRight(false).
			BorderForeground(lipgloss.Color("#9E9E9E"))
	}

	categoryPanel := categoryStyle.Render(params.CategoryListView)

	// Table panel with focus indicator
	tableStyle := lipgloss.NewStyle().
		Width(tableWidth).
		Height(availableHeight - 4)

	// Add border to show focus
	if params.StorageFocus == StorageFocusTable {
		tableStyle = tableStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(params.FocusedColor))
	} else {
		tableStyle = tableStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9E9E9E"))
	}

	tablePanel := tableStyle.Render(params.TableView)

	content := lipgloss.JoinHorizontal(lipgloss.Top, categoryPanel, tablePanel)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Storage"),
		searchBar,
		content,
	)
}

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
