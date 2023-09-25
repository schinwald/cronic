package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/schinwald/cronic/internal/styles"
)

type InputModel struct {
	title     string
	textInput textinput.Model
	err       error
}

func MakeInputModel() InputModel {
	ti := textinput.New()
	ti.Prompt = ""
	ti.PlaceholderStyle.Foreground(styles.DimmedForegroundColor)

	return InputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m InputModel) View() string {
	var view strings.Builder
	view.WriteString(m.textInput.View())
	return view.String()
}

func (m *InputModel) Focus() {
	m.textInput.Focus()
}

func (m *InputModel) Blur() {
	m.textInput.Blur()
}

func (m *InputModel) Placeholder(value string) {
	m.textInput.Placeholder = value
}

func (m *InputModel) Width(value int) {
	m.textInput.Width = value
}

func (m *InputModel) CharLimit(value int) {
	m.textInput.CharLimit = value
}

func (m *InputModel) Value() string {
	return m.textInput.Value()
}
