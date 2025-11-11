package shared

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModalType int

const (
	ModalNone ModalType = iota
	ModalHelp
	ModalSaveSearch
	ModalLoadSearch
)

type ModalView struct {
	active ModalType
	input  textinput.Model
}

// NewModalView creates and initializes a new ModalView component
func NewModalView() ModalView {
	input := textinput.New()
	input.Placeholder = "Enter name..."
	input.CharLimit = 50

	return ModalView{
		active: ModalNone,
		input:  input,
	}
}

// Update handles modal view-specific updates
func (m *ModalView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

// GetInputView returns the rendered input view
func (m *ModalView) GetInputView() string {
	return m.input.View()
}

// GetInputValue returns the current input value
func (m *ModalView) GetInputValue() string {
	return m.input.Value()
}

// ResetInput clears the input value
func (m *ModalView) ResetInput() {
	m.input.Reset()
}

// FocusInput focuses the input
func (m *ModalView) FocusInput() tea.Cmd {
	return m.input.Focus()
}

// BlurInput blurs the input
func (m *ModalView) BlurInput() {
	m.input.Blur()
}

// GetActive returns the currently active modal type
func (m *ModalView) GetActive() ModalType {
	return m.active
}

// SetActive sets the active modal type
func (m *ModalView) SetActive(modalType ModalType) {
	m.active = modalType
}

// IsActive returns whether any modal is active
func (m *ModalView) IsActive() bool {
	return m.active != ModalNone
}

// Close closes the current modal
func (m *ModalView) Close() {
	m.active = ModalNone
}
