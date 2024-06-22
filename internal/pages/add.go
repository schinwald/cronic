package pages

import (
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
	name = iota
	description
	command
	user
	flow
	ai
	schedule
	confirmation
)

type AddModel struct {
	Page
	state                int
	focus                int
	width                int
	height               int
	nameInput            components.InputModel
	descriptionInput     components.InputModel
	commandInput         components.InputModel
	userList             components.ListModel
	flowList             components.ListModel
	confirmationList     components.ListModel
	aiPanel              *layouts.AIModel
	schedulePanel        layouts.ScheduleModel
	legendPanel          *layouts.LegendModel
	nextOccurrencePanels []layouts.NextOccurrenceModel
	cronjob              *utils.CronJob
	err                  error
}

func OnFocus(l *layouts.LegendModel) func(int) error {
	return l.SetFocus
}

func MakeAddModel() *AddModel {
	nameInput := components.MakeInputModel()
	nameInput.Prompt("> ")
	nameInput.Placeholder("Dishes Reminder")
	nameInput.Focus()

	descriptionInput := components.MakeInputModel()
	descriptionInput.Prompt("> ")
	descriptionInput.Placeholder("Reminder to clean the dishes")

	commandInput := components.MakeInputModel()
	commandInput.Prompt("> ")
	commandInput.Placeholder("./dishes")

	userList := components.MakeListModel(utils.GetAllUsers())

	flowList := components.MakeListModel([]string{"ai (recommended)", "traditional"})

	confirmationList := components.MakeListModel([]string{"yes", "no"})

	// Pointer used here
	legendPanel := layouts.MakeLegendModel()

	schedulePanel := layouts.MakeScheduleModel()

	// Why does this work now since I am using a pointer above
	// What is the difference between struct{}, &(struct{}) and &struct{}
	schedulePanel.OnFocus(legendPanel.SetFocus)

	aiPanel := layouts.MakeAIModel()

	nextOccurrencePanels := []layouts.NextOccurrenceModel{
		layouts.MakeNextOccurrenceModel("First Occurrence"),
		layouts.MakeNextOccurrenceModel("Second Occurrence"),
		layouts.MakeNextOccurrenceModel("Third Occurrence"),
		layouts.MakeNextOccurrenceModel("Fourth Occurrence"),
		layouts.MakeNextOccurrenceModel("Fifth Occurrence"),
	}

	cronjob := utils.MakeCronJob()

	return &AddModel{
		state:                name,
		focus:                name,
		nameInput:            nameInput,
		descriptionInput:     descriptionInput,
		commandInput:         commandInput,
		userList:             userList,
		flowList:             flowList,
		confirmationList:     confirmationList,
		legendPanel:          legendPanel,
		schedulePanel:        schedulePanel,
		aiPanel:              aiPanel,
		nextOccurrencePanels: nextOccurrencePanels,
		cronjob:              cronjob,
	}
}

func (m AddModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.nameInput.Init())
	cmds = append(cmds, m.descriptionInput.Init())
	cmds = append(cmds, m.commandInput.Init())
	cmds = append(cmds, m.userList.Init())
	cmds = append(cmds, m.flowList.Init())
	cmds = append(cmds, m.confirmationList.Init())
	cmds = append(cmds, m.legendPanel.Init())
	cmds = append(cmds, m.schedulePanel.Init())
	cmds = append(cmds, m.aiPanel.Init())

	for _, nextOccurrencePanel := range m.nextOccurrencePanels {
		cmds = append(cmds, nextOccurrencePanel.Init())
	}

	return tea.Batch(cmds...)
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
			if m.state == m.focus {
				m.SubmitStep()
			} else {
				m.NextFocus()
			}
		case tea.KeyShiftTab:
			m.PreviousFocus()
		case tea.KeyTab:
			m.NextFocus()
		}
	case error:
		m.err = msg
		return m, nil
	}

	switch m.state {
	case name:
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)

		m.cronjob.Name(m.nameInput.Value())
	case description:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		cmds = append(cmds, cmd)

		m.cronjob.Description(m.descriptionInput.Value())
	case command:
		m.commandInput, cmd = m.commandInput.Update(msg)
		cmds = append(cmds, cmd)

		m.cronjob.Command(m.commandInput.Value())
	case user:
		m.userList, cmd = m.userList.Update(msg)
		cmds = append(cmds, cmd)

		m.cronjob.User(m.userList.Value())
	case flow:
		m.flowList, cmd = m.flowList.Update(msg)
		cmds = append(cmds, cmd)
	case ai:
		m.aiPanel, cmd = m.aiPanel.Update(msg)
		cmds = append(cmds, cmd)

		err = m.cronjob.Expression(m.aiPanel.CronExpression())
		if err != nil {
			for i, _ := range m.cronjob.Next {
				m.nextOccurrencePanels[i].NextOccurrence(time.Time{})
			}
			break
		}

		for i, next := range m.cronjob.Next {
			m.nextOccurrencePanels[i].NextOccurrence(next)
		}
	case schedule:
		m.schedulePanel, cmd = m.schedulePanel.Update(msg)
		cmds = append(cmds, cmd)

		m.legendPanel, cmd = m.legendPanel.Update(msg)
		cmds = append(cmds, cmd)

		err = m.cronjob.Expression(m.schedulePanel.CronExpression())
		if err != nil {
			m.schedulePanel.CronExplanation("")
			for i, _ := range m.cronjob.Next {
				m.nextOccurrencePanels[i].NextOccurrence(time.Time{})
			}
			break
		}

		m.schedulePanel.CronExplanation(m.cronjob.HumanReadable)
		for i, next := range m.cronjob.Next {
			m.nextOccurrencePanels[i].NextOccurrence(next)
		}
	case confirmation:
		m.confirmationList, cmd = m.confirmationList.Update(msg)
		cmds = append(cmds, cmd)
	}

	for i, _ := range m.nextOccurrencePanels {
		m.nextOccurrencePanels[i], cmd = m.nextOccurrencePanels[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m AddModel) View() string {
	var view strings.Builder

	inquiries := make(map[int]bool)
	inquiries[name] = true
	inquiries[description] = true
	inquiries[command] = true
	inquiries[user] = true
	inquiries[flow] = true

	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	if inquiries[m.state] {
		view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("What name would you give this cronjob?"),
			m.nameInput.View(),
			"",
		))
		if m.state == name {
			return view.String()
		}

		view.WriteRune('\n')

		view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("How would you describe the cronjob?"),
			m.descriptionInput.View(),
			"",
		))
		if m.state == description {
			return view.String()
		}

		view.WriteRune('\n')

		view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("What is the command?"),
			m.commandInput.View(),
			"",
		))
		if m.state == command {
			return view.String()
		}

		view.WriteRune('\n')

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
			titleStyle.Render("How would you like to create the cronjob?"),
			m.flowList.View(),
			"",
		))
		if m.state == flow {
			return view.String()
		}
	}

	nextOccurrencePanels := ""
	for _, nextOccurrencePanel := range m.nextOccurrencePanels {
		nextOccurrencePanel.Size(m.width, 4)

		nextOccurrencePanels = lipgloss.JoinHorizontal(lipgloss.Left,
			nextOccurrencePanels,
			nextOccurrencePanel.View(),
		)
	}

	if m.flowList.Value() == "ai (recommended)" {
		m.aiPanel.Size(67, 0)
		view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
			m.aiPanel.View(),
			nextOccurrencePanels,
		))
		if m.state == ai {
			return view.String()
		}
	} else if m.flowList.Value() == "traditional" {
		m.schedulePanel.Size(67, 0)
		m.legendPanel.Size(67, 0)
		view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
			m.schedulePanel.View(),
			m.legendPanel.View(),
			nextOccurrencePanels,
		))
		if m.state == schedule {
			return view.String()
		}
	}

	view.WriteRune('\n')

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Would you like to save the cronjob?"),
		m.confirmationList.View(),
	))
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
	case name:
		m.SetFocus(description)
		m.state = description
		return true
	case description:
		m.SetFocus(command)
		m.state = command
		return true
	case command:
		m.SetFocus(user)
		m.state = user
		return true
	case user:
		m.SetFocus(flow)
		m.state = flow
		return true
	case flow:
		if m.flowList.Value() == "ai (recommended)" {
			m.SetFocus(ai)
			m.state = ai
		} else if m.flowList.Value() == "traditional" {
			m.SetFocus(schedule)
			m.state = schedule
		}
		return true
	case ai:
		m.SetFocus(confirmation)
		m.state = confirmation
		return true
	case schedule:
		isNextFocused := m.schedulePanel.NextFocus()
		if !isNextFocused {
			m.SetFocus(confirmation)
			m.state = confirmation
		}
		return true
	case confirmation:
		if m.confirmationList.Value() == "yes" {
			err = m.cronjob.WriteCronJob()
			if err != nil {
				log.Fatal(err)
			}

			os.Exit(0)
		}

		if m.confirmationList.Value() == "no" {
			if m.flowList.Value() == "ai (recommended)" {
				m.SetFocus(ai)
				m.state = ai
			} else if m.flowList.Value() == "traditional" {
				m.SetFocus(schedule)
				m.state = schedule
			}
		}
	}

	return false
}

func (m *AddModel) Blur() error {
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.commandInput.Blur()
	m.userList.Blur()
	m.flowList.Blur()
	m.aiPanel.Blur()
	m.schedulePanel.Blur()
	m.confirmationList.Blur()

	return nil
}

func (m *AddModel) PreviousFocus() bool {
	switch m.focus {
	case name:
		return true
	case description:
		m.SetFocus(name)
		return true
	case command:
		m.SetFocus(description)
		return true
	case user:
		m.SetFocus(command)
		return true
	case flow:
		m.SetFocus(user)
		return true
	case ai:
		m.SetFocus(flow)
		m.state = flow
		return true
	case schedule:
		isPreviousFocused := m.schedulePanel.PreviousFocus()
		if !isPreviousFocused {
			m.SetFocus(flow)
			m.state = flow
		}
		return true
	case confirmation:
		if m.flowList.Value() == "ai (recommended)" {
			m.SetFocus(ai)
			m.state = ai
		} else if m.flowList.Value() == "traditional" {
			m.SetFocus(schedule)
			m.state = schedule
		}
		return true
	}

	return false
}

func (m *AddModel) NextFocus() bool {
	switch m.focus {
	case name:
		m.SetFocus(description)
		return true
	case description:
		m.SetFocus(command)
		return true
	case command:
		m.SetFocus(user)
		return true
	case user:
		m.SetFocus(flow)
		return true
	case flow:
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
	case name:
		m.nameInput.Focus()
		return nil
	case description:
		m.descriptionInput.Focus()
		return nil
	case command:
		m.commandInput.Focus()
		return nil
	case user:
		m.userList.Focus()
		return nil
	case flow:
		m.flowList.Focus()
		return nil
	case ai:
		m.aiPanel.Focus()
		return nil
	case schedule:
		m.schedulePanel.Focus()
		return nil
	case confirmation:
		m.confirmationList.Focus()
		return nil
	}

	return nil
}
