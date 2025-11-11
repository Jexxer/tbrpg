package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		m.leftTabsList.SetHeight(ws.LeftPanel.Height - ws.BorderOffset)
		m.detailsList.SetHeight(ws.DetailsPanel.Height - ws.BorderOffset)

		// Update activity viewport size
		m.activityViewport.Width = ws.ActivityPanel.Width
		m.activityViewport.Height = ws.ActivityPanel.Height - ws.BorderOffset

		// Update storage components
		m.storageCategoryList.SetHeight(4)
		m.storageTable.SetHeight(ws.MainPanel.Height - 10)

		return m, nil

	case tea.MouseMsg:
		// Handle mouse wheel scrolling in activity log
		if msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown {
			// Check if mouse is over activity log area (7 rows above command line)
			if msg.Y >= m.Height-10 && msg.Y < m.Height-3 {
				var cmd tea.Cmd
				m.activityViewport, cmd = m.activityViewport.Update(msg)
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
				m.commandMode = true
				m.commandInput.Focus()
				return m, textinput.Blink
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
		if m.activeModal != ModalNone {
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
		if m.commandMode {
			switch msg.String() {
			case "esc":
				m.commandMode = false
				m.commandInput.Blur()
				m.FocusedView = FocusGameView
				return m, nil
			case "enter":
				// Execute command (implement later)
				cmdText := m.commandInput.Value()
				styledText := lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
				m.AddLogEntry("Command", "Ran: ", styledText.Render(cmdText))

				m.commandInput.Reset()
				m.commandMode = false
				m.commandInput.Blur()
				m.FocusedView = FocusGameView
				return m, nil
			default:
				m.commandInput, cmd = m.commandInput.Update(msg)
				return m, cmd
			}
		}

		// Storage search mode handling
		if m.storageSearchActive && m.FocusedView == FocusGameView && m.ActiveTab == 1 {
			switch msg.String() {
			case "esc":
				m.storageSearchActive = false
				m.storageSearchInput.Blur()
				return m, nil
			case "enter":
				m.storageSearchActive = false
				m.storageSearchInput.Blur()
				m.AddLogEntry("Storage", "Search completed", "")
				return m, nil
			default:
				m.storageSearchInput, cmd = m.storageSearchInput.Update(msg)
				m.updateStorageTable()
				return m, cmd
			}
		}

		// Global keybindings
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case ":":
			m.FocusedView = FocusCommandLine
			m.commandMode = true
			m.commandInput.Focus()
			return m, textinput.Blink

		case "?":
			// Open help modal when in storage view
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 {
				m.activeModal = ModalHelp
				return m, nil
			}

		case "tab":
			m.FocusedView = (m.FocusedView + 1) % 5

		case "enter", " ":
			if m.FocusedView == FocusLeftTabs {
				m.ActiveTab = m.leftTabsList.Index()
				m.FocusedView = FocusGameView

				// TESTING: Add log entry when switching tabs
				tabName := []string{"Storage", "Navigation", "Equipment", "Gathering", "Processing", "Crafting", "Quests"}[m.ActiveTab]
				m.AddLogEntry("Navigation", "Switched to "+tabName, "")
			}
		case "/":
			// Activate search in storage view
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 {
				m.storageSearchActive = true
				m.storageSearchInput.Focus()
				return m, textinput.Blink
			}

		case "S":
			// Save current search
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 && m.storageSearchInput.Value() != "" {
				m.activeModal = ModalSaveSearch
				m.modalInput.Focus()
				return m, textinput.Blink
			}

		case "O":
			// Load saved search
			if m.FocusedView == FocusGameView && m.ActiveTab == 1 {
				m.activeModal = ModalLoadSearch
				return m, nil
			}
		}

		// Delegate to focused component
		switch m.FocusedView {
		case FocusLeftTabs:
			m.leftTabsList, cmd = m.leftTabsList.Update(msg)
			cmds = append(cmds, cmd)
		case FocusDetails:
			m.detailsList, cmd = m.detailsList.Update(msg)
			cmds = append(cmds, cmd)

		case FocusGameView:
			// Storage view specific handling
			if m.ActiveTab == 1 {
				// Don't handle keys if search is active
				if m.storageSearchActive {
					break
				}

				switch msg.String() {
				case "left", "h":
					// Switch focus to category list
					m.storageFocus = StorageFocusCategory
					m.storageTable.Blur() // Add this

				case "right", "l":
					// Switch focus to table
					m.storageFocus = StorageFocusTable
					m.storageTable.Focus() // Add this

				case "up", "k", "down", "j":
					// Route to focused component
					if m.storageFocus == StorageFocusCategory {
						m.storageCategoryList, cmd = m.storageCategoryList.Update(msg)

						// Update filter when category changes
						selectedItem := m.storageCategoryList.SelectedItem()
						if item, ok := selectedItem.(listItem); ok {
							if m.GameState.SelectedCategory != item.title {
								m.GameState.SetCategory(item.title)
								m.updateStorageTable()
								m.AddLogEntry("Storage", "Category: "+item.title, "")
							}
						}
						cmds = append(cmds, cmd)
					} else {
						m.storageTable, cmd = m.storageTable.Update(msg)
						cmds = append(cmds, cmd)
					}

				case "enter":
					// Handle selection
					if m.storageFocus == StorageFocusTable {
						// TODO: Show item details modal
						m.AddLogEntry("Storage", "Selected item", "")
					}

				default:
					// Pass other keys to table if focused
					if m.storageFocus == StorageFocusTable {
						m.storageTable, cmd = m.storageTable.Update(msg)
						cmds = append(cmds, cmd)
					}
				}
			}

		case FocusActivityLog:
			m.activityViewport, cmd = m.activityViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Always update command input if in command mode to keep cursor blinking
	if m.commandMode {
		m.commandInput, cmd = m.commandInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Handle modal input
func (m Model) handleModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.activeModal = ModalNone
		m.modalInput.Blur()
		return m, nil

	case "?":
		if m.activeModal == ModalHelp {
			m.activeModal = ModalNone
		}
		return m, nil

	case "enter":
		switch m.activeModal {
		case ModalSaveSearch:
			// Save the search
			searchName := m.modalInput.Value()
			searchQuery := m.storageSearchInput.Value()

			m.GameState.SaveSearch(searchName, searchQuery)

			m.AddLogEntry("Storage", "Saved search: "+searchName, "")
			m.modalInput.Reset()
			m.modalInput.Blur()
			m.activeModal = ModalNone
			return m, nil
		}

	default:
		// Delegate to modal input if applicable
		if m.activeModal == ModalSaveSearch {
			var cmd tea.Cmd
			m.modalInput, cmd = m.modalInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
