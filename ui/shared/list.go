package shared

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jexxer/tbrpg/ui/styles"
)

// ListItem is a simple list item implementation
type ListItem struct {
	TitleText string
}

func (i ListItem) Title() string       { return i.TitleText }
func (i ListItem) FilterValue() string { return i.TitleText }

// CompactDelegate is a minimal list delegate for compact display
type CompactDelegate struct{}

func (d CompactDelegate) Height() int                               { return 1 }
func (d CompactDelegate) Spacing() int                              { return 0 }
func (d CompactDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CompactDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(ListItem)
	if !ok {
		return
	}

	str := i.TitleText

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
