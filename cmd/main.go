package main

import (
	"log"
	"os"

	"github.com/schinwald/cronic/internal/pages"
	"github.com/schinwald/cronic/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func main() {
	var err error

	// Load cronjobs from the cronjob file
	utils.LoadCronJobs()

	// Create the audit log if it doesn't exist
	_, err = os.OpenFile("audit.log", os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

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
	close  bool
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
	return m.pages[m.state].Init()
}

func (m windowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.width, m.height, m.err = term.GetSize(0)

	if m.close {
		return m, tea.Batch(tea.ClearScreen, tea.Quit)
	}

	if m.err != nil {
		return m, nil
	}

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.close = true
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
	}

	return m, tea.Batch(cmds...)
}

func (m windowModel) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)

	if m.close {
		return ""
	}

	m.pages[m.state].Size(
		m.width-lipgloss.Width(style.Render("")),
		m.height-lipgloss.Height(style.Render("")),
	)

	return style.Render(m.pages[m.state].View())
}
