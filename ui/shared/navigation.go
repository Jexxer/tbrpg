package shared

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jexxer/tbrpg/ui/styles"
)

type NavigationView struct {
	tabsList list.Model
}

// NewNavigationView creates and initializes a new NavigationView component
func NewNavigationView() NavigationView {
	tabsItems := []list.Item{
		ListItem{TitleText: "Navigation"},
		ListItem{TitleText: "Storage"},
		ListItem{TitleText: "Equipment"},
		ListItem{TitleText: "Gathering"},
		ListItem{TitleText: "Processing"},
		ListItem{TitleText: "Crafting"},
		ListItem{TitleText: "Quests"},
	}

	tabsList := list.New(tabsItems, CompactDelegate{}, 15, 10)
	tabsList.SetShowHelp(false)
	tabsList.SetShowStatusBar(false)
	tabsList.SetFilteringEnabled(false)
	tabsList.SetShowTitle(false)

	return NavigationView{
		tabsList: tabsList,
	}
}

// Update handles navigation-specific updates
func (n *NavigationView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	n.tabsList, cmd = n.tabsList.Update(msg)
	return cmd
}

// View renders the navigation view
func (n *NavigationView) View() string {
	return n.tabsList.View()
}

// GetSelectedIndex returns the currently selected tab index
func (n *NavigationView) GetSelectedIndex() int {
	return n.tabsList.Index()
}

// UpdateSize updates the component sizes based on window dimensions
func (n *NavigationView) UpdateSize(width, height int) {
	ws := styles.GetWindowSizes(width, height)
	n.tabsList.SetHeight(ws.LeftPanel.Height - ws.BorderOffset)
}
