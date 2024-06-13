package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schinwald/cronic/internal/styles"
)

var (
	focusedPromptStyle = lipgloss.NewStyle().Foreground(styles.ForegroundColor)
	blurredPromptStyle = lipgloss.NewStyle().Foreground(styles.TernaryColor)
)

type InputModel struct {
	textInput   textinput.Model
	prompt      string
	promptStyle lipgloss.Style
	err         error
}

func MakeInputModel() InputModel {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = ""

	return InputModel{
		textInput:   ti,
		promptStyle: focusedPromptStyle,
		err:         nil,
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

	view.WriteString(m.promptStyle.Render(m.prompt))
	view.WriteString(m.textInput.View())

	return view.String()
}

func (m *InputModel) Focus() {
	m.textInput.Focus()
	m.promptStyle = focusedPromptStyle
}

func (m *InputModel) Blur() {
	m.textInput.Blur()
	m.promptStyle = blurredPromptStyle
}

func (m *InputModel) Prompt(value string) {
	m.prompt = value
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
