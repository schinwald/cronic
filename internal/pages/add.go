package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/layouts"
)

type AddModel struct {
	fileInput           components.InputModel
	descriptionInput    components.InputModel
	schedulePanel       layouts.ScheduleModel
	legendPanel         layouts.LegendModel
	nextOccurrencePanel layouts.NextOccurrenceModel
	err                 error
}

func MakeAddModel() AddModel {
	return AddModel{
		fileInput: components.MakeInputModel("File", "./script", 150, 20),
		descriptionInput: components.MakeInputModel("Description", "Executes a command", 150, 20),
		schedulePanel: layouts.MakeScheduleModel(),
		legendPanel: layouts.MakeLegendModel(),
		nextOccurrencePanel: layouts.MakeNextOccurrenceModel(),
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	return m, tea.Batch(cmds...)
}

func (m AddModel) View() string {
	var view strings.Builder

	var main string
	main = lipgloss.JoinHorizontal(lipgloss.Top, m.schedulePanel.View(), m.legendPanel.View())
	main = lipgloss.JoinVertical(lipgloss.Left, main, m.nextOccurrencePanel.View())

	view.WriteString(m.fileInput.View())
	view.WriteString(m.descriptionInput.View())
	view.WriteRune('\n')
	view.WriteRune('\n')
	view.WriteString(main)

	return view.String()
}
