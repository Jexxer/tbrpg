package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/ui/shared"
	"github.com/jexxer/tbrpg/ui/styles"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	ws := styles.GetWindowSizes(m.Width, m.Height)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Update component sizes
		m.navigation.UpdateSize(msg.Width, msg.Height)
		m.storage.UpdateSize(msg.Width, msg.Height)
		m.activity.UpdateSize(msg.Width, msg.Height)
		m.detailsList.SetHeight(ws.DetailsPanel.Height - ws.BorderOffset)

		return m, nil

	case tea.MouseMsg:
		// Handle mouse wheel scrolling in activity log
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			// Check if mouse is over activity log area (7 rows above command line)
			if msg.Y >= m.Height-10 && msg.Y < m.Height-3 {
				cmd := m.activity.Update(msg)
				return m, cmd
			}
		}

		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			x, y := msg.X, msg.Y

			// Top bar (rows 0-2)
			if y >= 0 && y <= ws.TopPanel.Height-1 {
				return m, nil
			}

			// Activity log (7 rows above command line)
			if y >= m.Height-ws.CommandPanel.Height-ws.ActivityPanel.Height && y < m.Height-ws.CommandPanel.Height {
				m.FocusedView = FocusActivityLog
				return m, nil
			}

			// Command line (bottom 3 rows)
			if y >= m.Height-ws.CommandPanel.Height {
				m.FocusedView = FocusCommandLine
				cmd := m.command.Activate()
				return m, cmd
			}

			// Middle section (between top bar and activity log)
			// if below top bar and above activity box
			if y >= ws.TopPanel.Height && y < m.Height-ws.CommandPanel.Height-ws.ActivityPanel.Height {
				// if x less than or equal to left pannel width
				if x <= ws.LeftPanel.Width-1 {
					m.FocusedView = FocusLeftTabs

					// if we are in x bounds for right side
				} else if x >= m.Width-ws.CharacterInfoPanel.Width {
					// if we are below the characterInfoPanel (we check above in activity log conditionals)
					if y >= ws.TopPanel.Height+ws.CharacterInfoPanel.Height {
						m.FocusedView = FocusDetails
					}
				} else {
					m.FocusedView = FocusGameView
				}
			}
		}

	case tea.KeyMsg:
		// Modal handling (highest priority)
		if m.modal.IsActive() {
			return m.handleModalInput(msg)
		}

		// handle navigating panes via vim keys
		switch m.FocusedView {
		case FocusLeftTabs:
			switch msg.String() {
			case "L":
				m.FocusedView = FocusGameView
			case "J":
				m.FocusedView = FocusActivityLog
			}
		case FocusGameView:
			switch msg.String() {
			case "H":
				m.FocusedView = FocusLeftTabs
			case "L":
				m.FocusedView = FocusDetails
			case "J":
				m.FocusedView = FocusActivityLog
			}
		case FocusDetails:
			switch msg.String() {
			case "H":
				m.FocusedView = FocusGameView
			case "J":
				m.FocusedView = FocusActivityLog
			}
		case FocusActivityLog:
			switch msg.String() {
			case "K":
				m.FocusedView = FocusGameView
			case "J":
				m.FocusedView = FocusCommandLine
			}
		case FocusCommandLine:
			switch msg.String() {
			case "K":
				m.FocusedView = FocusActivityLog
			}
		}

		// Command mode handling
		if m.command.IsActive() {
			switch msg.String() {
			case "esc":
				m.command.Deactivate()
				m.FocusedView = FocusGameView
				return m, nil
			case "enter":
				// Execute command (implement later)
				cmdText := m.command.GetValue()
				styledText := lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
				m.AddLogEntry("Command", "Ran: ", styledText.Render(cmdText))

				m.command.Reset()
				m.command.Deactivate()
				m.FocusedView = FocusGameView
				return m, nil
			default:
				cmd = m.command.Update(msg)
				return m, cmd
			}
		}

		// Storage search mode handling (delegated to storage component in storage update section below)

		// Global keybindings
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case ":":
			m.FocusedView = FocusCommandLine
			cmd = m.command.Activate()
			return m, cmd

		case "?":
			// Open help modal when in storage view
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 {
				m.modal.SetActive(shared.ModalHelp)
				return m, nil
			}

		case "tab":
			m.FocusedView = (m.FocusedView + 1) % 5

		case "enter", " ":
			if m.FocusedView == FocusLeftTabs {
				m.ActiveTab = m.navigation.GetSelectedIndex()
				m.FocusedView = FocusGameView

				// TESTING: Add log entry when switching tabs
				tabName := []string{"Navigation", "Storage", "Equipment", "Gathering", "Processing", "Crafting", "Quests"}[m.ActiveTab]
				m.AddLogEntry("Navigation", "Switched to "+tabName, "")
			}
		case "/":
			// Handled by storage component

		case "S":
			// Save current search (handled by storage component but modal needs to be set)
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 && !m.storage.IsSearchActive() {
				m.modal.SetActive(shared.ModalSaveSearch)
				cmd = m.modal.FocusInput()
				return m, cmd
			}

		case "O":
			// Load saved search
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 {
				m.modal.SetActive(shared.ModalLoadSearch)
				return m, nil
			}
		}

		// Delegate to focused component
		switch m.FocusedView {
		case FocusLeftTabs:
			cmd = m.navigation.Update(msg)
			cmds = append(cmds, cmd)
		case FocusDetails:
			m.detailsList, cmd = m.detailsList.Update(msg)
			cmds = append(cmds, cmd)

		case FocusGameView:
			// Storage view specific handling
			if m.ActiveTab == 1 {
				cmd = m.storage.Update(msg, m.GameState, m.AddLogEntry)
				cmds = append(cmds, cmd)
			}

		case FocusActivityLog:
			cmd = m.activity.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Always update command input if in command mode to keep cursor blinking
	if m.command.IsActive() {
		cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Handle modal input
func (m Model) handleModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modal.Close()
		m.modal.BlurInput()
		return m, nil

	case "?":
		if m.modal.GetActive() == shared.ModalHelp {
			m.modal.Close()
		}
		return m, nil

	case "enter":
		switch m.modal.GetActive() {
		case shared.ModalSaveSearch:
			// Save the search
			searchName := m.modal.GetInputValue()
			// Note: We'll need to get the search query from storage component
			// For now, using a placeholder - this will be fixed when we refactor further

			m.GameState.SaveSearch(searchName, "placeholder")

			m.AddLogEntry("Storage", "Saved search: "+searchName, "")
			m.modal.ResetInput()
			m.modal.BlurInput()
			m.modal.Close()
			return m, nil
		}

	default:
		// Delegate to modal input if applicable
		if m.modal.GetActive() == shared.ModalSaveSearch {
			cmd := m.modal.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
