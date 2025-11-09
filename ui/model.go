package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
)

type FocusedView int


type (
	ModalType    int
	StorageFocus int
)


const (
	StorageFocusCategory StorageFocus = iota
	StorageFocusTable
)

const (
	ModalNone ModalType = iota
	ModalHelp
	ModalSaveSearch
	ModalLoadSearch
)

const (
	FocusLeftTabs FocusedView = iota
	FocusGameView
	FocusQuickActions
	FocusActivityLog
	FocusCommandLine
)

type Model struct {
	Width        int
	Height       int
	FocusedView  FocusedView
	ActiveTab    int // Which tab's content is showing in game view
	FocusedColor string

	// Game state
	GameState *game.State

	// Bubbles components
	leftTabsList     list.Model
	quickActionsList list.Model
	resourcesTable   table.Model
	commandInput     textinput.Model
	commandMode      bool

	// Activity Log
	activityViewport viewport.Model

	// Storage system
	storageSearchInput  textinput.Model
	storageCategoryList list.Model
	storageTable        table.Model
	storageSearchActive bool
	storageFocus        StorageFocus

	// Modal system
	activeModal ModalType
	modalInput  textinput.Model
}

type listItem struct {
	title string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.title }

type compactDelegate struct{}

func (d compactDelegate) Height() int                               { return 1 }
func (d compactDelegate) Spacing() int                              { return 0 }
func (d compactDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d compactDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	str := i.title

	// Highlight selected item
	if index == m.Index() {
		str = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#27F576")).
			Bold(true).
			Render("> " + str)
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

func InitialModel() Model {
	leftTabsItems := []list.Item{
		listItem{title: "Navigation"},
		listItem{title: "Storage"},
		listItem{title: "Equipment"},
		listItem{title: "Gathering"},
		listItem{title: "Processing"},
		listItem{title: "Crafting"},
		listItem{title: "Quests"},
	}

	leftTabsList := list.New(leftTabsItems, compactDelegate{}, 15, 10)
	leftTabsList.SetShowHelp(false)
	leftTabsList.SetShowStatusBar(false)
	leftTabsList.SetFilteringEnabled(false)
	leftTabsList.SetShowTitle(false)

	// Setup quick actions list
	quickActionsItems := []list.Item{
		listItem{title: "[A]ttack"},
		listItem{title: "[G]ather"},
		listItem{title: "[C]raft"},
		listItem{title: "[I]nventory"},
	}

	quickActionsList := list.New(quickActionsItems, compactDelegate{}, 20, 5)
	quickActionsList.SetShowHelp(false)
	quickActionsList.SetShowStatusBar(false)
	quickActionsList.SetFilteringEnabled(false)
	quickActionsList.SetShowTitle(false)

	// Setup resources table
	columns := []table.Column{
		{Title: "Resource", Width: 15},
		{Title: "Quantity", Width: 10},
	}

	rows := []table.Row{
		{"Wood", "150"},
		{"Iron Ore", "45"},
		{"Coal", "23"},
		{"Fish", "12"},
		{"Stone", "89"},
		{"Gold", "1,234g"},
	}

	resourcesTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	// Setup command input
	commandInput := textinput.New()
	// commandInput.Placeholder = "Enter command..." // TODO: bugged rn, only shows first char of placeholder
	commandInput.CharLimit = 100

	// Initialize game state
	gameState := game.NewState()

	// Setup storage search input
	storageSearchInput := textinput.New()
	storageSearchInput.Placeholder = "Search items... (? for help)"
	storageSearchInput.CharLimit = 100
	storageSearchInput.Width = 50

	// Setup storage category list
	categories := game.GetCategories()
	categoryItems := make([]list.Item, len(categories))
	for i, cat := range categories {
		categoryItems[i] = listItem{title: cat.Name}
	}

	storageCategoryList := list.New(categoryItems, compactDelegate{}, 20, 8)
	storageCategoryList.SetShowHelp(false)
	storageCategoryList.SetShowStatusBar(false)
	storageCategoryList.SetFilteringEnabled(false)
	storageCategoryList.SetShowTitle(false)

	// Setup storage table
	storageColumns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Qty", Width: 15},
		{Title: "Value (g)", Width: 8},
		// {Title: "Tags", Width: 25},
	}

	storageTable := table.New(
		table.WithColumns(storageColumns),
		table.WithRows([]table.Row{}), // Will populate with filtered results
		table.WithFocused(false),
		table.WithHeight(10),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#27F576")).
		Bold(false)
	storageTable.SetStyles(tableStyles)

	// Setup modal input
	modalInput := textinput.New()
	modalInput.Placeholder = "Enter name..."
	modalInput.CharLimit = 50

	// Setup activity viewport
	activityViewport := viewport.New(78, 5)
	activityViewport.SetContent("Activity log initialized...")

	m := Model{
		Width:               80,
		Height:              24,
		FocusedView:         FocusGameView,
		FocusedColor:        "#27F576",
		ActiveTab:           0,
		GameState:           gameState,
		leftTabsList:        leftTabsList,
		quickActionsList:    quickActionsList,
		resourcesTable:      resourcesTable,
		commandInput:        commandInput,
		commandMode:         false,
		activityViewport:    activityViewport,
		storageSearchInput:  storageSearchInput,
		storageCategoryList: storageCategoryList,
		storageTable:        storageTable,
		storageSearchActive: false,
		activeModal:         ModalNone,
		modalInput:          modalInput,
		storageFocus:        StorageFocusCategory,
	}

	m.updateActivityViewport()
	m.activityViewport.GotoBottom()

	// Update sotrage table with all items initiallly
	m.updateStorageTable()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// AddLogEntry entry to activity log
func (m *Model) AddLogEntry(category, action, details string) {
	m.GameState.ActivityLog.AddEntry(category, action, details)

	// Update viewport content
	m.updateActivityViewport()

	// Auto-scroll to bottom if not focused
	if m.FocusedView != FocusActivityLog {
		m.activityViewport.GotoBottom()
	}
}

func (m *Model) formatActivityLog() string {
	var content strings.Builder

	// Include ALL entries from game state
	for _, entry := range m.GameState.ActivityLog.GetEntries() {
		timestamp := entry.Timestamp.Format("15:04")

		color := getCategoryColor(entry.Category)

		categoryStyle := lipgloss.NewStyle().
			Foreground(color).
			Bold(true)

		// Fixed-width category (12 characters)
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

// Update viewport content
func (m *Model) updateActivityViewport() {
	content := m.formatActivityLog()
	m.activityViewport.SetContent(content)
}


// Update storage table with filtered items
func (m *Model) updateStorageTable() {
	searchTerm := m.storageSearchInput.Value()
	filtered := m.GameState.GetFilteredItems(searchTerm)

	rows := make([]table.Row, len(filtered))
	for i, item := range filtered {

		// Right-align Qty and Value columns
		qtyStr := fmt.Sprintf("%d", item.Quantity)
		valueStr := fmt.Sprintf("%d", item.Value)

		rows[i] = table.Row{
			item.Name,
			fmt.Sprintf("%8s", qtyStr),   // Right-align with padding
			fmt.Sprintf("%8s", valueStr), // Right-align with padding
		}
	}

	m.storageTable.SetRows(rows)
}
