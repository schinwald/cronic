package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)
type inputModel struct {
	textInput textinput.Model
	err error 
}

func initialModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "Cron Job Notation"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputModel{
		textInput: ti,
		err: nil, 
	}
}

func (im inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (im inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (im inputModel) View() string {
	return fmt.Sprintf(
		"Input Cron Job Notation:\n\n%s\n\n%s",
		im.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func inputCronNotation() {
	
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}