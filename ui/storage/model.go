package storage

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/game"
	"github.com/jexxer/tbrpg/ui/styles"
)

type Focus int

const (
	FocusCategory Focus = iota
	FocusTable
)

type View struct {
	searchInput  textinput.Model
	categoryList list.Model
	table        table.Model
	searchActive bool
	focus        Focus
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

	// Highlight selected item
	if index == m.Index() {
		str = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.FocusedColor)).
			Bold(true).
			Render("> " + str)
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
}

// New creates and initializes a new storage View
func New() View {
	// Setup storage search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search items... (? for help)"
	searchInput.CharLimit = 100
	searchInput.Width = 25

	// Setup storage category list
	categories := game.GetCategories()
	categoryItems := make([]list.Item, len(categories))
	for i, cat := range categories {
		categoryItems[i] = listItem{title: cat.Name}
	}

	categoryList := list.New(categoryItems, compactDelegate{}, 20, 4)
	categoryList.SetShowPagination(false)
	categoryList.SetShowHelp(false)
	categoryList.SetShowStatusBar(false)
	categoryList.SetFilteringEnabled(false)
	categoryList.SetShowTitle(false)

	// Setup storage table
	storageColumns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Qty", Width: 15},
		{Title: "Value (g)", Width: 8},
	}

	storageTable := table.New(
		table.WithColumns(storageColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(false),
		table.WithHeight(20),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.BorderColor)).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color(styles.BlackColor)).
		Background(lipgloss.Color(styles.FocusedColor)).
		Bold(false)
	storageTable.SetStyles(tableStyles)

	return View{
		searchInput:  searchInput,
		categoryList: categoryList,
		table:        storageTable,
		searchActive: false,
		focus:        FocusCategory,
	}
}

// IsSearchActive returns whether search mode is active
func (v *View) IsSearchActive() bool {
	return v.searchActive
}
