package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/schinwald/cronic/internal/pages"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const (
	add = iota
	list
)

type windowModel struct {
	state int
	pages []tea.Model
	err error
}

func initialModel() windowModel {
	pages := []tea.Model{
		pages.MakeAddModel(),
	}

	return windowModel{
		state: add,
		pages: pages,
	}
}

func (m windowModel) Init() tea.Cmd {
	return nil
}

func (m windowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil	
	}

	switch m.state {
	// Handle all actions from initial state
	case add:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return m, tea.Quit
			}
		case error:
			m.err = msg
			return m, nil	
		}
	// Handle all actions from list state
	case list:
		break
	}

	return m, tea.Batch(cmds...)
}

func (m windowModel) View() string {
	return m.pages[m.state].View()
}
