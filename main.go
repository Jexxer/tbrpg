package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int
}

func initialModel() model {
	return model{
		width:  80,
		height: 24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	// Calculate dimensions
	leftWidth := 15
	rightWidth := 25
	centerWidth := m.width - leftWidth - rightWidth - 6
	contentHeight := m.height - 8

	// Top bar
	topBar := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.width - 2).
		Render(" Location: Varrock | Current Action: Woodcutting")

	// Left tabs
	leftTabs := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(leftWidth).
		Height(contentHeight).
		Render(m.renderLeftTabs())

	// Game view
	gameView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(centerWidth).
		Height(contentHeight).
		Render(m.renderGameView())

	// Character info (top right)
	charInfoHeight := contentHeight / 2
	charInfo := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(rightWidth).
		Height(charInfoHeight).
		Render(m.renderCharacterInfo())

	quickActionsHeight := contentHeight - charInfoHeight
	// Quick actions (bottom right)
	quickActions := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(rightWidth).
		Height(quickActionsHeight - 2).
		Render(m.renderQuickActions())

	// Combine right side vertically
	rightSide := lipgloss.JoinVertical(lipgloss.Left, charInfo, quickActions)

	// Combine middle section horizontally
	middleSection := lipgloss.JoinHorizontal(lipgloss.Top, leftTabs, gameView, rightSide)

	// Command line
	commandLine := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.width - 2).
		Render(" > commands to play (press 'q' to quit)")

	// Combine everything vertically
	return lipgloss.JoinVertical(lipgloss.Left, topBar, middleSection, commandLine)
}

func (m model) renderLeftTabs() string {
	tabs := []string{"Storage", "Navigation", "Equipment", "Gathering", "Processing", "Crafting", "Quests"}
	content := ""
	for _, tab := range tabs {
		content += tab + "\n"
	}
	return content
}

func (m model) renderGameView() string {
	return "Game View\n\nYou are in Varrock\nA bustling city full of opportunities"
}

func (m model) renderCharacterInfo() string {
	return "Character Info\n\nName: Adventurer\nHealth: 50/50\nMana: 30/30\nLevel: 15\nGold: 1234g"
}

func (m model) renderQuickActions() string {
	return "Quick Actions\n\n[A]ttack\n[G]ather\n[C]raft\n[I]nventory"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
