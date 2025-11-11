// Package ui provides the terminal user interface components
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/ui/shared"
	"github.com/jexxer/tbrpg/ui/storage"
	"github.com/jexxer/tbrpg/ui/styles"
)

func (m Model) View() string {
	baseView := m.renderBaseView()

	// Overlay modal if active
	if m.modal.IsActive() {
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

	leftTabs := leftTabsStyle.Render(m.navigation.View())

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

	activityLog := activityStyle.Render(m.activity.View())

	// Command line
	commandLineStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(windowStyles.CommandPanel.Width - windowStyles.BorderOffset)

	if m.FocusedView == FocusCommandLine {
		commandLineStyle = commandLineStyle.BorderForeground(lipgloss.Color(styles.FocusedColor))
	} else {
		commandLineStyle = commandLineStyle.BorderForeground(lipgloss.Color(styles.UnfocusedColor))
	}

	commandLine := commandLineStyle.Render(m.command.View())

	return lipgloss.JoinVertical(lipgloss.Left, topBar, middleSection, activityLog, commandLine)
}

func (m Model) renderModalOverlay(base string) string {
	var modalContent string

	switch m.modal.GetActive() {
	case shared.ModalHelp:
		modalContent = m.renderHelpModal()
	case shared.ModalSaveSearch:
		modalContent = m.renderSaveSearchModal()
	case shared.ModalLoadSearch:
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
	return storage.RenderHelpModal()
}

func (m Model) renderSaveSearchModal() string {
	// TODO: Get search query from storage component
	return storage.RenderSaveSearchModal(
		"placeholder",
		m.modal.GetInputView(),
	)
}

func (m Model) renderLoadSearchModal() string {
	return storage.RenderLoadSearchModal(m.GameState.SavedSearches)
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
	return m.storage.View(m.Width, m.Height)
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
