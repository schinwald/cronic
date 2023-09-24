package pages

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/layouts"
	"github.com/schinwald/cronic/internal/styles"
)

const (
	file = iota
	description
	schedule
	legend
	nextOccurrence
)

type AddModel struct {
	Page
	state               int
	focus               int
	width               int
	height              int
	fileInput           components.InputModel
	descriptionInput    components.InputModel
	schedulePanel       layouts.ScheduleModel
	legendPanel         *layouts.LegendModel
	nextOccurrencePanel layouts.NextOccurrenceModel
	err                 error
}

func OnFocus(l *layouts.LegendModel) func (int) error {
	return l.SetFocus
}

func MakeAddModel() *AddModel {
	fileInput := components.MakeInputModel()
	fileInput.Focus()
	fileInput.Placeholder("./dishes")

	descriptionInput := components.MakeInputModel()
	descriptionInput.Placeholder("Reminder to clean the dishes")

	schedulePanel := layouts.MakeScheduleModel()
	// Pointer used here
	legendPanel := layouts.MakeLegendModel()
	// Why does this work now since I am using a pointer above
	// What is the difference between struct{}, &(struct{}) and &struct{}
	schedulePanel.OnFocus(legendPanel.SetFocus)

	nextOccurrencePanel := layouts.MakeNextOccurrenceModel()

	return &AddModel{
		state:               file,
		focus:               file,
		fileInput:           fileInput,
		descriptionInput:    descriptionInput,
		schedulePanel:       schedulePanel,
		legendPanel:         legendPanel,
		nextOccurrencePanel: nextOccurrencePanel,
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (Page, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.SubmitStep()
		case tea.KeyShiftTab:
			m.PreviousFocus()
		case tea.KeyTab:
			m.NextFocus()
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.fileInput, cmd = m.fileInput.Update(msg)
	cmds = append(cmds, cmd)

	m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	cmds = append(cmds, cmd)

	m.schedulePanel, cmd = m.schedulePanel.Update(msg)
	cmds = append(cmds, cmd)

	m.legendPanel, cmd = m.legendPanel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AddModel) View() string {
	var view strings.Builder

	gap := lipgloss.NewStyle().Padding(1, 1).Render("")

	legendWidth := 70
	m.schedulePanel.Size(m.width-(lipgloss.Width(gap)/2 + legendWidth), 15)
	m.legendPanel.Size(legendWidth - lipgloss.Width(gap)/2, 15)
	m.nextOccurrencePanel.Size(int(m.width*1.0), 4)

	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	view.WriteString(titleStyle.Render("Program: "))
	view.WriteString(m.fileInput.View())
	if m.state == file {
		return view.String()
	}

	view.WriteRune('\n')
	view.WriteString(titleStyle.Render("Description: "))
	view.WriteString(m.descriptionInput.View())
	if m.state == description {
		return view.String()
	}

	var main string
	main = lipgloss.JoinHorizontal(lipgloss.Top, m.schedulePanel.View(), gap, m.legendPanel.View())
	main = lipgloss.JoinVertical(lipgloss.Left, main, m.nextOccurrencePanel.View())
	view.WriteRune('\n')
	view.WriteRune('\n')
	view.WriteRune('\n')
	view.WriteString(main)
	return view.String()
}

func (m *AddModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *AddModel) SubmitStep() error {
	switch m.state {
	case file:
		m.state = description
		m.NextFocus()
		return nil
	case description:
		m.state = schedule
		m.NextFocus()
		return nil
	case schedule:
		m.NextFocus()
		return nil
	}

	return nil
}

func (m *AddModel) Blur() error {
	m.fileInput.Blur()
	m.descriptionInput.Blur()
	m.schedulePanel.Blur()

	return nil
}

func (m *AddModel) PreviousFocus() error {
	switch m.focus {
	case file:
		return errors.New("done")
	case description:
		m.SetFocus(file)
		return nil
	case schedule:
		err := m.schedulePanel.PreviousFocus()
		if err != nil {
			m.SetFocus(description)
		}
		return err
	}

	return nil
}

func (m *AddModel) NextFocus() error {
	switch m.focus {
	case file:
		m.SetFocus(description)
		return nil
	case description:
		m.SetFocus(schedule)
		return nil
	case schedule:
		err := m.schedulePanel.NextFocus()
		return err
	}

	return nil
}

func (m *AddModel) SetFocus(focus int) error {
	m.focus = focus

	m.Blur()

	switch focus {
	case file:
		m.fileInput.Focus()
		return nil
	case description:
		m.descriptionInput.Focus()
		return nil
	case schedule:
		m.schedulePanel.Focus()
		return nil
	}

	return nil
}
