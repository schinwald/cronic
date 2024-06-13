package pages

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/layouts"
	"github.com/schinwald/cronic/internal/styles"
	"github.com/schinwald/cronic/internal/utils"
)

const (
	user = iota
	command
	name
	description
	schedule
	confirmation
)

type AddModel struct {
	Page
	state               int
	focus               int
	width               int
	height              int
	userList            components.ListModel
	commandInput        components.InputModel
	nameInput           components.InputModel
	descriptionInput    components.InputModel
	confirmationInput   components.InputModel
	schedulePanel       layouts.ScheduleModel
	legendPanel         *layouts.LegendModel
	nextOccurrencePanel layouts.NextOccurrenceModel
	cronjob             *utils.CronJob
	err                 error
}

func OnFocus(l *layouts.LegendModel) func(int) error {
	return l.SetFocus
}

func MakeAddModel() *AddModel {
	userList := components.MakeListModel(utils.GetAllUsers())

	commandInput := components.MakeInputModel()
	commandInput.Placeholder("./dishes")

	nameInput := components.MakeInputModel()
	nameInput.Placeholder("Dishes Reminder")

	descriptionInput := components.MakeInputModel()
	descriptionInput.Placeholder("Reminder to clean the dishes")

	confirmationInput := components.MakeInputModel()
	confirmationInput.Placeholder("Yes")

	schedulePanel := layouts.MakeScheduleModel()
	// Pointer used here
	legendPanel := layouts.MakeLegendModel()
	// Why does this work now since I am using a pointer above
	// What is the difference between struct{}, &(struct{}) and &struct{}
	schedulePanel.OnFocus(legendPanel.SetFocus)

	nextOccurrencePanel := layouts.MakeNextOccurrenceModel()

	cronjob := utils.MakeCronJob()

	return &AddModel{
		state:               user,
		focus:               user,
		userList:            userList,
		commandInput:        commandInput,
		nameInput:           nameInput,
		descriptionInput:    descriptionInput,
		confirmationInput:   confirmationInput,
		schedulePanel:       schedulePanel,
		legendPanel:         legendPanel,
		nextOccurrencePanel: nextOccurrencePanel,
		cronjob:             cronjob,
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (Page, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var err error

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

	m.userList, cmd = m.userList.Update(msg)
	cmds = append(cmds, cmd)

	m.commandInput, cmd = m.commandInput.Update(msg)
	cmds = append(cmds, cmd)

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	cmds = append(cmds, cmd)

	m.schedulePanel, cmd = m.schedulePanel.Update(msg)
	cmds = append(cmds, cmd)

	m.confirmationInput, cmd = m.confirmationInput.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case user:
		m.cronjob.User(m.userList.Value())
	case command:
		m.cronjob.Command(m.commandInput.Value())
	case name:
		m.cronjob.Name(m.nameInput.Value())
	case description:
		m.cronjob.Description(m.descriptionInput.Value())
	case schedule:
		err = m.cronjob.Expression(m.schedulePanel.CronExpression())
		if err != nil {
			m.schedulePanel.CronExplanation("")
			m.nextOccurrencePanel.NextOccurrence(time.Time{})
			break
		}

		m.schedulePanel.CronExplanation(m.cronjob.HumanReadable)
		m.nextOccurrencePanel.NextOccurrence(m.cronjob.Next)
	case confirmation:
	}

	m.legendPanel, cmd = m.legendPanel.Update(msg)
	cmds = append(cmds, cmd)

	m.nextOccurrencePanel, cmd = m.nextOccurrencePanel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AddModel) View() string {
	var view strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Who do you want to run this command?"),
		m.userList.View(),
		"",
	))
	if m.state == user {
		return view.String()
	}

	view.WriteRune('\n')

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("What is the command?"),
		fmt.Sprintf("> %s", m.commandInput.View()),
		"",
	))
	if m.state == command {
		return view.String()
	}

	view.WriteRune('\n')

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("What name would you give this cronjob?"),
		fmt.Sprintf("> %s", m.nameInput.View()),
		"",
	))
	if m.state == name {
		return view.String()
	}

	view.WriteRune('\n')

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("How would you describe the cronjob?"),
		fmt.Sprintf("> %s", m.descriptionInput.View()),
		"",
	))
	if m.state == description {
		return view.String()
	}

	m.schedulePanel.Size(67, 0)
	m.legendPanel.Size(67, 0)
	m.nextOccurrencePanel.Size(m.width, 4)
	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		"",
		"",
		m.nextOccurrencePanel.View(),
		m.schedulePanel.View(),
		m.legendPanel.View(),
	))
	if m.state == schedule {
		return view.String()
	}

	view.WriteRune('\n')

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Would you like to save the cronjob (Y/n)?"),
		fmt.Sprintf("> %s", m.confirmationInput.View()),
	))
	view.WriteRune('\n')
	if m.state == confirmation {
		return view.String()
	}

	return view.String()
}

func (m *AddModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *AddModel) SubmitStep() bool {
	var err error

	switch m.state {
	case user:
		m.SetFocus(command)
		m.state = command
		return true
	case command:
		m.SetFocus(name)
		m.state = name
		return true
	case name:
		m.SetFocus(description)
		m.state = description
		return true
	case description:
		m.SetFocus(schedule)
		m.state = schedule
		return true
	case schedule:
		isNextFocused := m.schedulePanel.NextFocus()

		if !isNextFocused {
			m.SetFocus(confirmation)
			m.state = confirmation
		}

		return true
	case confirmation:
		err = m.cronjob.WriteCronJob()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	return false
}

func (m *AddModel) Blur() error {
	m.userList.Blur()
	m.commandInput.Blur()
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.schedulePanel.Blur()
	m.confirmationInput.Blur()

	return nil
}

func (m *AddModel) PreviousFocus() bool {
	switch m.focus {
	case user:
		return false
	case command:
		m.SetFocus(user)
		return true
	case name:
		m.SetFocus(command)
		return true
	case description:
		m.SetFocus(name)
		return true
	case schedule:
		isPreviousFocused := m.schedulePanel.PreviousFocus()
		if !isPreviousFocused {
			m.SetFocus(description)
		}
		return true
	case confirmation:
		m.SetFocus(schedule)
		return true
	}

	return false
}

func (m *AddModel) NextFocus() bool {
	switch m.focus {
	case user:
		m.SetFocus(command)
		return true
	case command:
		m.SetFocus(name)
		return true
	case name:
		m.SetFocus(description)
		return true
	case description:
		m.SetFocus(schedule)
		return true
	case schedule:
		isNextFocused := m.schedulePanel.NextFocus()

		if !isNextFocused {
			m.SetFocus(confirmation)
		}

		return isNextFocused
	case confirmation:
		return false
	}

	return false
}

func (m *AddModel) SetFocus(focus int) error {
	m.focus = focus

	m.Blur()

	switch focus {
	case user:
		m.userList.Focus()
		return nil
	case command:
		m.commandInput.Focus()
		return nil
	case name:
		m.nameInput.Focus()
		return nil
	case description:
		m.descriptionInput.Focus()
		return nil
	case schedule:
		m.schedulePanel.Focus()
		return nil
	case confirmation:
		m.confirmationInput.Focus()
		return nil
	}

	return nil
}
