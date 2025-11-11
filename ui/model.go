package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jexxer/tbrpg/game"
	"github.com/jexxer/tbrpg/ui/shared"
	"github.com/jexxer/tbrpg/ui/storage"
)

type FocusedView int

const (
	FocusLeftTabs FocusedView = iota
	FocusGameView
	FocusDetails
	FocusActivityLog
	FocusCommandLine
)

type Model struct {
	Width       int
	Height      int
	FocusedView FocusedView
	ActiveTab   int // Which tab's content is showing in game view

	// Game state
	GameState *game.State

	// View components
	navigation shared.NavigationView
	storage    storage.View
	activity   shared.ActivityView
	command    shared.CommandView
	modal      shared.ModalView

	// Legacy components (to be refactored later)
	detailsList    list.Model
	resourcesTable table.Model
}

type listItem struct {
	title string
}

func (i listItem) Title() string       { return i.title }
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

	// Highlight selected item (placeholder for legacy components)
	if index == m.Index() {
		str = "> " + str
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

func InitialModel() Model {
	// Setup details list (legacy component - to be refactored)
	detailsItems := []list.Item{
		listItem{title: "[A]ttack"},
		listItem{title: "[G]ather"},
		listItem{title: "[C]raft"},
		listItem{title: "[I]nventory"},
	}

	detailsList := list.New(detailsItems, compactDelegate{}, 20, 5)
	detailsList.SetShowHelp(false)
	detailsList.SetShowStatusBar(false)
	detailsList.SetFilteringEnabled(false)
	detailsList.SetShowTitle(false)

	// Setup resources table (legacy component - to be refactored)
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

	// Initialize game state
	gameState := game.NewState()

	// Initialize view components
	navigation := shared.NewNavigationView()
	storageView := storage.New()
	activity := shared.NewActivityView()
	command := shared.NewCommandView()
	modal := shared.NewModalView()

	m := Model{
		Width:          80,
		Height:         24,
		FocusedView:    FocusGameView,
		ActiveTab:      0,
		GameState:      gameState,
		navigation:     navigation,
		storage:        storageView,
		activity:       activity,
		command:        command,
		modal:          modal,
		detailsList:    detailsList,
		resourcesTable: resourcesTable,
	}

	// Initialize activity log
	m.activity.UpdateContent(gameState)
	m.activity.GotoBottom()

	// Update storage table with all items initially
	m.storage.UpdateTable(gameState)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// AddLogEntry adds an entry to the activity log
func (m *Model) AddLogEntry(category, action, details string) {
	m.GameState.ActivityLog.AddEntry(category, action, details)

	// Update activity view content
	m.activity.UpdateContent(m.GameState)

	// Auto-scroll to bottom if not focused
	if m.FocusedView != FocusActivityLog {
		m.activity.GotoBottom()
	}
}
