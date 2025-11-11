// Package shared
package shared

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
	"github.com/jexxer/tbrpg/ui/styles"
)

type ActivityView struct {
	viewport viewport.Model
}

// NewActivityView creates and initializes a new ActivityView component
func NewActivityView() ActivityView {
	vp := viewport.New(78, 5)
	vp.SetContent("Activity log initialized...")

	return ActivityView{
		viewport: vp,
	}
}

// Update handles activity view-specific updates
func (a *ActivityView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return cmd
}

// View renders the activity view
func (a *ActivityView) View() string {
	return a.viewport.View()
}

// UpdateContent updates the viewport content with the activity log
func (a *ActivityView) UpdateContent(gameState *game.State) {
	content := a.formatActivityLog(gameState)
	a.viewport.SetContent(content)
}

// GotoBottom scrolls to the bottom of the viewport
func (a *ActivityView) GotoBottom() {
	a.viewport.GotoBottom()
}

// formatActivityLog formats the activity log entries for display
func (a *ActivityView) formatActivityLog(gameState *game.State) string {
	var content strings.Builder

	for _, entry := range gameState.ActivityLog.GetEntries() {
		timestamp := entry.Timestamp.Format("15:04")
		color := getCategoryColor(entry.Category)

		categoryStyle := lipgloss.NewStyle().
			Foreground(color).
			Bold(true)

		paddedCategory := fmt.Sprintf("%-11s", entry.Category)

		line := fmt.Sprintf("[%s] %s  %s",
			timestamp,
			categoryStyle.Render(paddedCategory),
			entry.Action,
		)

		if entry.Details != "" {
			line += " " + entry.Details
		}

		content.WriteString("\n" + line)
	}

	return content.String()
}

// getCategoryColor returns the color for a category
func getCategoryColor(category string) lipgloss.Color {
	categoryColors := map[string]lipgloss.Color{
		"Combat":      "196",
		"Woodcutting": "34",
		"Fishing":     "33",
		"Market":      "226",
		"Navigation":  "205",
		"System":      "240",
		"Command":     "240",
		"Storage":     "33",
	}

	if color, ok := categoryColors[category]; ok {
		return color
	}
	return "255"
}

// UpdateSize updates the component sizes based on window dimensions
func (a *ActivityView) UpdateSize(width, height int) {
	ws := styles.GetWindowSizes(width, height)
	a.viewport.Width = ws.ActivityPanel.Width
	a.viewport.Height = ws.ActivityPanel.Height - ws.BorderOffset
}
