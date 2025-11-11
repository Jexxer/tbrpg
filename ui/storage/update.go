package storage

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jexxer/tbrpg/game"
	"github.com/jexxer/tbrpg/ui/styles"
)

// UpdateTable updates the storage table with filtered items from game state
func (v *View) UpdateTable(gameState *game.State) {
	searchTerm := v.searchInput.Value()
	filtered := gameState.GetFilteredItems(searchTerm)

	rows := make([]table.Row, len(filtered))
	for i, item := range filtered {
		qtyStr := fmt.Sprintf("%d", item.Quantity)
		valueStr := fmt.Sprintf("%d", item.Value)

		rows[i] = table.Row{
			item.Name,
			fmt.Sprintf("%8s", qtyStr),
			fmt.Sprintf("%8s", valueStr),
		}
	}

	v.table.SetRows(rows)
}

// Update handles storage-specific updates
func (v *View) Update(msg tea.Msg, gameState *game.State, onLog func(category, action, details string)) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode
		if v.searchActive {
			switch msg.String() {
			case "esc":
				v.searchActive = false
				v.searchInput.Blur()
				return nil
			case "enter":
				v.searchActive = false
				v.searchInput.Blur()
				onLog("Storage", "Search completed", "")
				return nil
			default:
				v.searchInput, cmd = v.searchInput.Update(msg)
				v.UpdateTable(gameState)
				return cmd
			}
		}

		// Handle normal keyboard navigation
		switch msg.String() {
		case "/":
			v.searchActive = true
			v.searchInput.Focus()
			return textinput.Blink

		case "left", "h":
			v.focus = FocusCategory
			v.table.Blur()

		case "right", "l":
			v.focus = FocusTable
			v.table.Focus()

		case "up", "k", "down", "j":
			if v.focus == FocusCategory {
				v.categoryList, cmd = v.categoryList.Update(msg)

				// Update filter when category changes
				selectedItem := v.categoryList.SelectedItem()
				if item, ok := selectedItem.(listItem); ok {
					if gameState.SelectedCategory != item.title {
						gameState.SetCategory(item.title)
						v.UpdateTable(gameState)
						onLog("Storage", "Category: "+item.title, "")
					}
				}
				cmds = append(cmds, cmd)
			} else {
				v.table, cmd = v.table.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "enter":
			if v.focus == FocusTable {
				onLog("Storage", "Selected item", "")
			}

		default:
			if v.focus == FocusTable {
				v.table, cmd = v.table.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

// UpdateSize updates the component sizes based on window dimensions
func (v *View) UpdateSize(width, height int) {
	ws := styles.GetWindowSizes(width, height)
	v.categoryList.SetHeight(4)
	v.table.SetHeight(ws.MainPanel.Height - 10)
}
