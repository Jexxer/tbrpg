// Package styles is for holding all relevant common styles
package styles

import "github.com/charmbracelet/lipgloss"

// Layout dimensions
const (
	LeftPanelWidth    = 15
	RightPanelWidth   = 25
	PanelSpacing      = 6
	ContentHeightSub  = 15
	ActivityViewportW = 78
	ActivityViewportH = 5
)

// Colors
const (
	FocusedColor   = "#27F576"
	UnfocusedColor = "#9E9E9E"
	WhiteColor     = "#FFFFFF"
	BlackColor     = "#000000"
	BorderColor    = "240"
)

// CategoryColors for activity log
var CategoryColors = map[string]lipgloss.Color{
	"Combat":      "196",
	"Woodcutting": "34",
	"Fishing":     "33",
	"Market":      "226",
	"Navigation":  "205",
	"System":      "240",
	"Command":     "240",
	"Storage":     "33",
}

// GetCategoryColor returns the color for a category, or white if not found
func GetCategoryColor(category string) lipgloss.Color {
	if color, ok := CategoryColors[category]; ok {
		return color
	}
	return "255" // Default white
}

// FocusedBorderStyle return common styles
func FocusedBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(FocusedColor))
}

func UnfocusedBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(UnfocusedColor))
}

func SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(FocusedColor)).
		Bold(true)
}

// GetTableSelectedStyle table styles
func GetTableSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BlackColor)).
		Background(lipgloss.Color(FocusedColor)).
		Bold(false)
}

func GetTableHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(BorderColor)).
		BorderBottom(true).
		Bold(false)
}
