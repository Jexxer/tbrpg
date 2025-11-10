// Package styles is for holding all relevant common styles
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

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

type styles struct {
	Height int
	Width  int
}

type storage struct {
	Categories styles
	Table      styles
	Search     styles
}

type WindowStyles struct {
	TopPanel           styles
	LeftPanel          styles
	MainPanel          styles
	CharacterInfoPanel styles
	DetailsPanel       styles
	ActivityPanel      styles
	CommandPanel       styles
	BorderOffset       int
	Storage            storage
}

func GetWindowSizes(terminalWidth int, terminalHeight int) *WindowStyles {
	topPanelHeight := 3
	topPanelWidth := terminalWidth

	leftPanelWidth := 24

	characterInfoPanelWidth := 24
	characterInfoPanelHeight := 9

	detailsPanelWidth := characterInfoPanelWidth

	activityPanelHeight := 7
	activityPanelWidth := terminalWidth

	commandPanelHeight := 3
	commandPanelWidth := terminalWidth

	// Caclulate remaining
	mainPanelWidth := terminalWidth - leftPanelWidth - characterInfoPanelWidth
	mainPanelHeight := terminalHeight - topPanelHeight - activityPanelHeight - commandPanelHeight

	leftPanelHeight := terminalHeight - topPanelHeight - activityPanelHeight - commandPanelHeight

	detailsPanelHeight := terminalHeight - topPanelHeight - characterInfoPanelHeight - activityPanelHeight - commandPanelHeight

	// Storage View Styles
	categoryWidth := 20
	categoryHeight := 10

	searchInputHeight := 3
	searchInputWidth := mainPanelWidth

	tableHeight := mainPanelHeight - searchInputHeight
	tableWidth := mainPanelWidth - categoryWidth

	storageStyles := storage{
		Categories: styles{Width: categoryWidth, Height: categoryHeight},
		Table:      styles{Width: tableWidth, Height: tableHeight},
		Search:     styles{Width: searchInputWidth, Height: searchInputHeight},
	}

	return &WindowStyles{
		TopPanel:           styles{Width: topPanelWidth, Height: topPanelHeight},
		LeftPanel:          styles{Width: leftPanelWidth, Height: leftPanelHeight},
		MainPanel:          styles{Width: mainPanelWidth, Height: mainPanelHeight},
		CharacterInfoPanel: styles{Width: characterInfoPanelWidth, Height: characterInfoPanelHeight},
		DetailsPanel:       styles{Width: detailsPanelWidth, Height: detailsPanelHeight},
		ActivityPanel:      styles{Width: activityPanelWidth, Height: activityPanelHeight},
		CommandPanel:       styles{Width: commandPanelWidth, Height: commandPanelHeight},
		BorderOffset:       2,
		Storage:            storageStyles,
	}
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
