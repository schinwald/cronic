package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schinwald/cronic/internal/pages"
	"golang.org/x/term"
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
	state  int
	pages  []pages.Page
	width  int
	height int
	err    error
}

func initialModel() windowModel {
	pages := []pages.Page{
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
	var cmd tea.Cmd

	m.width, m.height, m.err = term.GetSize(0)

	if m.err != nil {
		return m, nil
	}

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
		m.pages[add], cmd = m.pages[add].Update(msg)
		cmds = append(cmds, cmd)
	// Handle all actions from list state
	case list:
		m.pages[list], cmd = m.pages[list].Update(msg)
		cmds = append(cmds, cmd)
		break
	}

	return m, tea.Batch(cmds...)
}

func (m windowModel) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)

	m.pages[m.state].Size(
		m.width - lipgloss.Width(style.Render("")), 
		m.height - lipgloss.Height(style.Render("")),
	)

	return style.Render(m.pages[m.state].View())
}
