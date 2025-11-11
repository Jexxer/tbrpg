package shared

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandView struct {
	input textinput.Model
	mode  bool
}

// NewCommandView creates and initializes a new CommandView component
func NewCommandView() CommandView {
	input := textinput.New()
	input.Placeholder = "Enter command..."
	input.CharLimit = 100
	input.Width = 25

	return CommandView{
		input: input,
		mode:  false,
	}
}

// Update handles command view-specific updates
func (c *CommandView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return cmd
}

// View renders the command view
func (c *CommandView) View() string {
	if c.mode {
		return c.input.View()
	}
	return " > Commands - ? for help - Press ':' to enter command mode"
}

// Activate enables command mode and focuses the input
func (c *CommandView) Activate() tea.Cmd {
	c.mode = true
	c.input.Focus()
	return textinput.Blink
}

// Deactivate disables command mode and blurs the input
func (c *CommandView) Deactivate() {
	c.mode = false
	c.input.Blur()
}

// Reset clears the input value
func (c *CommandView) Reset() {
	c.input.Reset()
}

// GetValue returns the current input value
func (c *CommandView) GetValue() string {
	return c.input.Value()
}

// IsActive returns whether command mode is active
func (c *CommandView) IsActive() bool {
	return c.mode
}

// Focus focuses the input
func (c *CommandView) Focus() tea.Cmd {
	return c.input.Focus()
}

// Blur blurs the input
func (c *CommandView) Blur() {
	c.input.Blur()
}
