package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/styles"
)

type InputModel struct {
	title string
	textInput textinput.Model
	err error 
}

func MakeInputModel(title string, placeholder string, charLimit int, width int) InputModel {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = width

	return InputModel{
		title: title,
		textInput: ti,
		err: nil, 
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
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

	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	view.WriteString(titleStyle.Render(m.title + ":"))
	view.WriteString(m.textInput.View())
	view.WriteRune('\n')

	return view.String()
}
