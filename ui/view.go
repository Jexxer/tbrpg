// Package ui provides the terminal user interface components
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/ui/styles"
	"github.com/jexxer/tbrpg/ui/views"
)

func (m Model) View() string {
	baseView := m.renderBaseView()

	// Overlay modal if active
	if m.activeModal != ModalNone {
		return m.renderModalOverlay(baseView)
	}

	return baseView
}

func (m Model) renderBaseView() string {
	// leftWidth := 15
	// rightWidth := 25
	// centerWidth := m.Width - leftWidth - rightWidth - 6
	// contentHeight := m.Height - 15
	windowStyles := styles.GetWindowSizes(m.Width, m.Height)

	// Top bar
	topBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.TopPanel.Width - windowStyles.BorderOffset).
		Align(lipgloss.Center).
		Render("The Town of Starting")

	// Left tabs - use list component
	leftTabsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.LeftPanel.Width - windowStyles.BorderOffset).
		Height(windowStyles.LeftPanel.Height - windowStyles.BorderOffset)

	if m.FocusedView == FocusLeftTabs {
		leftTabsStyle = leftTabsStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		leftTabsStyle = leftTabsStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	leftTabs := leftTabsStyle.Render(m.leftTabsList.View())

	// Game view
	gameViewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.MainPanel.Width - windowStyles.BorderOffset).
		Height(windowStyles.MainPanel.Height - windowStyles.BorderOffset)

	if m.FocusedView == FocusGameView {
		gameViewStyle = gameViewStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		gameViewStyle = gameViewStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	gameView := gameViewStyle.Render(m.renderGameView())

	// Character info
	charInfo := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.CharacterInfoPanel.Width - windowStyles.BorderOffset).
		Height(windowStyles.CharacterInfoPanel.Height - windowStyles.BorderOffset).
		BorderForeground(lipgloss.Color(styles.UnfocusedColor)).
		Render(m.renderCharacterInfo())

	// Details - use list component
	detailsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.DetailsPanel.Width - windowStyles.BorderOffset).
		Height(windowStyles.DetailsPanel.Height - windowStyles.BorderOffset)

	if m.FocusedView == FocusDetails {
		detailsStyle = detailsStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		detailsStyle = detailsStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	details := detailsStyle.Render(m.detailsList.View())

	rightSide := lipgloss.JoinVertical(lipgloss.Left, charInfo, details)
	middleSection := lipgloss.JoinHorizontal(lipgloss.Top, leftTabs, gameView, rightSide)

	// Activity log panel
	activityStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.ActivityPanel.Width - windowStyles.BorderOffset).
		Height(windowStyles.ActivityPanel.Height - windowStyles.BorderOffset)

	if m.FocusedView == FocusActivityLog {
		activityStyle = activityStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		activityStyle = activityStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	activityLog := activityStyle.Render(m.activityViewport.View())

	// Command line
	var commandLineContent string
	if m.commandMode {
		commandLineContent = m.commandInput.View()
	} else {
		commandLineContent = " > Commands - ? for help - Press ':' to enter command mode"
	}

	commandLineStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.CommandPanel.Width - windowStyles.BorderOffset)

	if m.FocusedView == FocusCommandLine {
		commandLineStyle = commandLineStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		commandLineStyle = commandLineStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	commandLine := commandLineStyle.Render(commandLineContent)

	return lipgloss.JoinVertical(lipgloss.Left, topBar, middleSection, activityLog, commandLine)
}

func (m Model) renderModalOverlay(base string) string {
	var modalContent string

	switch m.activeModal {
	case ModalHelp:
		modalContent = m.renderHelpModal()
	case ModalSaveSearch:
		modalContent = m.renderSaveSearchModal()
	case ModalLoadSearch:
		modalContent = m.renderLoadSearchModal()
	default:
		return base
	}

	// Create modal box
	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(60).
		Render(modalContent)

	// Overlay on base view using Place
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("236")),
	)
}

func (m Model) renderHelpModal() string {
	return views.RenderHelpModal()
}

func (m Model) renderSaveSearchModal() string {
	return views.RenderSaveSearchModal(
		m.storageSearchInput.Value(),
		m.modalInput.View(),
	)
}

func (m Model) renderLoadSearchModal() string {
	return views.RenderLoadSearchModal(m.GameState.SavedSearches)
}

func (m Model) renderGameView() string {
	switch m.ActiveTab {
	case 0: // Navigation
		return m.renderNavigationView()
	case 1: // Storage
		return m.renderStorageView()
	case 2: // Equipment
		return m.renderEquipmentView()
	case 3: // Gathering
		return m.renderGatheringView()
	case 4: // Processing
		return m.renderProcessingView()
	case 5: // Crafting
		return m.renderCraftingView()
	case 6: // Quests
		return m.renderQuestsView()
	default:
		return "Unknown view"
	}
}

func (m Model) renderStorageView() string {
	ws := styles.GetWindowSizes(m.Width, m.Height)

	// Calculate available dimensions
	availableWidth := ws.MainPanel.Width - ws.BorderOffset
	availableHeight := ws.MainPanel.Height - ws.BorderOffset
	searchBarHeight := 3

	// Update input width
	m.storageSearchInput.Width = availableWidth - 33

	// Update table dimensions
	qtyWidth := 4
	valueWidth := 6
	nameWidth := 10

	m.storageTable.Columns()[0].Width = nameWidth
	m.storageTable.Columns()[1].Width = qtyWidth
	m.storageTable.Columns()[2].Width = valueWidth
	m.storageTable.SetHeight(availableHeight - searchBarHeight - ws.BorderOffset)

	// Map StorageFocus to int for views package
	var focusInt int
	if m.storageFocus == StorageFocusCategory {
		focusInt = views.StorageFocusCategory
	} else {
		focusInt = views.StorageFocusTable
	}

	return views.RenderStorageView(views.StorageViewParams{
		Width:               m.Width,
		Height:              m.Height,
		FocusedColor:        styles.FocusedColor,
		StorageSearchActive: m.storageSearchActive,
		StorageFocus:        focusInt,
		SearchInputView:     m.storageSearchInput.View(),
		CategoryListView:    m.storageCategoryList.View(),
		TableView:           m.storageTable.View(),
		AvailableWidth:      availableWidth,
		AvailableHeight:     availableHeight,
	})
}

func (m Model) renderCharacterInfo() string {
	return "Character Info\n\nName: Adventurer\nHealth: 50/50\nMana: 30/30\nLevel: 15\nGold: 1234g"
}

// Keep your other render functions for now
func (m Model) renderNavigationView() string {
	return "Navigation View\n\nComing soon..."
}

func (m Model) renderEquipmentView() string {
	return "Equipment View\n\nComing soon..."
}

func (m Model) renderGatheringView() string {
	return "Gathering View\n\nComing soon..."
}

func (m Model) renderProcessingView() string {
	return "Processing View\n\nComing soon..."
}

func (m Model) renderCraftingView() string {
	return "Crafting View\n\nComing soon..."
}

func (m Model) renderQuestsView() string {
	return "Quests View\n\nComing soon..."
}
