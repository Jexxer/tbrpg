package storage

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
	"github.com/jexxer/tbrpg/ui/styles"
)

// View renders the storage view
func (v *View) View(width, height int) string {
	ws := styles.GetWindowSizes(width, height)

	// Calculate available dimensions
	availableWidth := ws.MainPanel.Width - ws.BorderOffset
	availableHeight := ws.MainPanel.Height - ws.BorderOffset
	searchBarHeight := 3

	// Update input width
	v.searchInput.Width = availableWidth - 33

	// Update table dimensions
	qtyWidth := 10
	valueWidth := 10
	nameWidth := availableWidth - ws.Storage.Categories.Width - qtyWidth - valueWidth - 9

	v.table.Columns()[0].Width = nameWidth
	v.table.Columns()[1].Width = qtyWidth
	v.table.Columns()[2].Width = valueWidth
	v.table.SetHeight(availableHeight - searchBarHeight - ws.BorderOffset)

	// Render the view
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	// Search bar
	searchBarStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Padding(0, 1)

	searchBar := searchBarStyle.Render("Search: " + v.searchInput.View() + " (press / to search)")

	// Layout: Category list on left, table on right
	tableWidth := ws.MainPanel.Width - ws.Storage.Categories.Width - (ws.BorderOffset * 2) - 1

	// Category panel with focus indicator
	categoryStyle := lipgloss.NewStyle().
		Width(ws.Storage.Categories.Width).
		Height(len(game.GetCategories()))

	// Add border to show focus
	if v.focus == FocusCategory {
		categoryStyle = categoryStyle.
			Border(lipgloss.RoundedBorder()).
			BorderRight(false).
			BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		categoryStyle = categoryStyle.
			Border(lipgloss.RoundedBorder()).
			BorderRight(false).
			BorderForeground(lipgloss.Color("#9E9E9E"))
	}

	categoryPanel := categoryStyle.Render(v.categoryList.View())

	// Table panel with focus indicator
	tableStyle := lipgloss.NewStyle().
		Width(tableWidth)
		// Height(ws.MainPanel.Height - ws.Storage.Search.Height - (ws.BorderOffset * 2))

	// Add border to show focus
	if v.focus == FocusTable {
		tableStyle = tableStyle.
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		tableStyle = tableStyle.
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#9E9E9E"))
	}

	tablePanel := tableStyle.Render(v.table.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, categoryPanel, tablePanel)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Storage"),
		searchBar,
		content,
	)
}
