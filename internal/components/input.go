package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type InputModel struct {
	textInput textinput.Model
	err error
}

func InitialModel() InputModel {
	ti := textinput.New()
	ti.Placeholder = "Cron Job Notation"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return InputModel{
		textInput: ti,
		err: nil,
	}
}

func (im InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (im InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return im, tea.Quit
		}
	case errMsg:
		im.err = msg
		return im, nil
	}

	im.textInput, cmd = im.textInput.Update(msg) 
	return im, cmd
}

func (im InputModel) View() string {
	return fmt.Sprintf(
		"Input Cron Job Notation:\n\n%s\n\n%s",
		im.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
