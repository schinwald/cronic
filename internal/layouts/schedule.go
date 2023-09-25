package layouts

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/styles"
)

type ScheduleModel struct {
	focus           int
	focusCallback   func(int) error
	width           int
	height          int
	minuteInput     components.InputModel
	hourInput       components.InputModel
	dayOfMonthInput components.InputModel
	monthInput      components.InputModel
	dayOfWeekInput  components.InputModel
	err             error
}

func defaultFocusCallback(focus int) error {
	return nil
}

func MakeScheduleModel() ScheduleModel {
	minuteInput := components.MakeInputModel()
	hourInput := components.MakeInputModel()
	dayOfMonthInput := components.MakeInputModel()
	monthInput := components.MakeInputModel()
	dayOfWeekInput := components.MakeInputModel()

	return ScheduleModel{
		focus:           minute,
		focusCallback:   defaultFocusCallback,
		minuteInput:     minuteInput,
		hourInput:       hourInput,
		dayOfMonthInput: dayOfMonthInput,
		monthInput:      monthInput,
		dayOfWeekInput:  dayOfWeekInput,
	}
}

func (m ScheduleModel) Init() tea.Cmd {
	return nil
}

func (m ScheduleModel) Update(msg tea.Msg) (ScheduleModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			switch m.focus {
			case minute:
				m.NextFocus()
				break
			case hour:
				m.NextFocus()
				break
			case dayOfMonth:
				m.NextFocus()
				break
			case month:
				m.NextFocus()
				break
			case dayOfWeek:
				break
			}
		case tea.KeyBackspace:
			switch m.focus {
			case minute:
				break
			case hour:
				if m.hourInput.Value() == "" {
					m.PreviousFocus()
				}
				break
			case dayOfMonth:
				if m.dayOfMonthInput.Value() == "" {
					m.PreviousFocus()
				}
				break
			case month:
				if m.monthInput.Value() == "" {
					m.PreviousFocus()
				}
				break
			case dayOfWeek:
				if m.dayOfWeekInput.Value() == "" {
					m.PreviousFocus()
				}
				break
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.minuteInput, cmd = m.minuteInput.Update(msg)
	cmds = append(cmds, cmd)

	m.hourInput, cmd = m.hourInput.Update(msg)
	cmds = append(cmds, cmd)

	m.dayOfMonthInput, cmd = m.dayOfMonthInput.Update(msg)
	cmds = append(cmds, cmd)

	m.monthInput, cmd = m.monthInput.Update(msg)
	cmds = append(cmds, cmd)

	m.dayOfWeekInput, cmd = m.dayOfWeekInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ScheduleModel) View() string {
	var view strings.Builder
	var content string

	paddingY, paddingX := 0, 0

	naturalTextBar := lipgloss.NewStyle().Italic(true).Render("Every minute")
	inputBar := createInputBar(m, 12)
	headerBar := createHeaderBar(m, 12)

	content = lipgloss.JoinVertical(lipgloss.Center, naturalTextBar, "", inputBar, headerBar)
	content = lipgloss.Place(m.width-2, m.height-2, lipgloss.Center, lipgloss.Center, content)

	view.WriteString(styles.PanelStyle("Schedule", content, m.width, m.height, paddingY, paddingX))

	return view.String()
}

func (m *ScheduleModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *ScheduleModel) Focus() error {
	m.SetFocus(m.focus)

	return nil
}

func (m *ScheduleModel) Blur() error {
	m.minuteInput.Blur()
	m.hourInput.Blur()
	m.dayOfMonthInput.Blur()
	m.monthInput.Blur()
	m.dayOfWeekInput.Blur()

	return nil
}

func (m *ScheduleModel) PreviousFocus() error {
	switch m.focus {
	case minute:
		return errors.New("done")
	case hour:
		m.SetFocus(minute)
		return nil
	case dayOfMonth:
		m.SetFocus(hour)
		return nil
	case month:
		m.SetFocus(dayOfMonth)
		return nil
	case dayOfWeek:
		m.SetFocus(month)
		return nil
	}

	return nil
}

func (m *ScheduleModel) NextFocus() error {
	switch m.focus {
	case minute:
		m.SetFocus(hour)
		return nil
	case hour:
		m.SetFocus(dayOfMonth)
		return nil
	case dayOfMonth:
		m.SetFocus(month)
		return nil
	case month:
		m.SetFocus(dayOfWeek)
		return nil
	case dayOfWeek:
		return errors.New("done")
	}

	return nil
}

func (m *ScheduleModel) OnFocus(callback func(int) error) error {
	m.focusCallback = callback

	return nil
}

func (m *ScheduleModel) SetFocus(focus int) error {
	m.focusCallback(focus)
	m.focus = focus

	m.Blur()

	switch focus {
	case minute:
		m.minuteInput.Focus()
		return nil
	case hour:
		m.hourInput.Focus()
		return nil
	case dayOfMonth:
		m.dayOfMonthInput.Focus()
		return nil
	case month:
		m.monthInput.Focus()
		return nil
	case dayOfWeek:
		m.dayOfWeekInput.Focus()
		return nil
	}

	return nil
}

func createInputBar(m ScheduleModel, width int) string {
	var temp string
	var inputs = make([]strings.Builder, 5)
	var inputStyles = make([]lipgloss.Style, 5)

	inputStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	inputs[minute].WriteString(m.minuteInput.View())
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, inputs[minute].String())
	inputs[minute].Reset()
	inputs[minute].WriteString(temp)
	inputStyles[minute] = inputStyle

	inputs[hour].WriteString(m.hourInput.View())
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, inputs[hour].String())
	inputs[hour].Reset()
	inputs[hour].WriteString(temp)
	inputStyles[hour] = inputStyle

	inputs[dayOfMonth].WriteString(m.dayOfMonthInput.View())
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, inputs[dayOfMonth].String())
	inputs[dayOfMonth].Reset()
	inputs[dayOfMonth].WriteString(temp)
	inputStyles[dayOfMonth] = inputStyle

	m.monthInput.CharLimit(width)
	inputs[month].WriteString(m.monthInput.View())
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, inputs[month].String())
	inputs[month].Reset()
	inputs[month].WriteString(temp)
	inputStyles[month] = inputStyle

	inputs[dayOfWeek].WriteString(m.dayOfWeekInput.View())
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, inputs[dayOfWeek].String())
	inputs[dayOfWeek].Reset()
	inputs[dayOfWeek].WriteString(temp)
	inputStyles[dayOfWeek] = inputStyle

	inputValues := lipgloss.JoinHorizontal(lipgloss.Top,
		inputStyles[minute].Render(inputs[minute].String()),
		inputStyles[hour].Render(inputs[hour].String()),
		inputStyles[dayOfMonth].Render(inputs[dayOfMonth].String()),
		inputStyles[month].Render(inputs[month].String()),
		inputStyles[dayOfWeek].Render(inputs[dayOfWeek].String()),
	)

	return styles.BorderStyle.Copy().Padding(0, 1).Render(inputValues)
}

func createHeaderBar(m ScheduleModel, width int) string {
	var temp string
	var headers = make([]strings.Builder, 5)
	var headerStyles = make([]lipgloss.Style, 5)

	headerStyle := lipgloss.NewStyle().Align(lipgloss.Center)
	inactiveHeaderStyle := lipgloss.NewStyle().Inherit(headerStyle).Faint(true)
	activeHeaderStyle := lipgloss.NewStyle().Inherit(headerStyle).Bold(true)

	headers[minute].WriteString("minute")
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, headers[minute].String())
	headers[minute].Reset()
	headers[minute].WriteString(temp)
	if m.minuteInput.Value() != "" || m.focus == minute {
		headerStyles[minute] = activeHeaderStyle
	} else {
		headerStyles[minute] = inactiveHeaderStyle
	}

	headers[hour].WriteString("hour")
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, headers[hour].String())
	headers[hour].Reset()
	headers[hour].WriteString(temp)
	if m.hourInput.Value() != "" || m.focus == hour {
		headerStyles[hour] = activeHeaderStyle
	} else {
		headerStyles[hour] = inactiveHeaderStyle
	}

	headers[dayOfMonth].WriteString("day")
	headers[dayOfMonth].WriteRune('\n')
	headers[dayOfMonth].WriteString("of month")
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, headers[dayOfMonth].String())
	headers[dayOfMonth].Reset()
	headers[dayOfMonth].WriteString(temp)
	if m.dayOfMonthInput.Value() != "" || m.focus == dayOfMonth {
		headerStyles[dayOfMonth] = activeHeaderStyle
	} else {
		headerStyles[dayOfMonth] = inactiveHeaderStyle
	}

	headers[month].WriteString("month")
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, headers[month].String())
	headers[month].Reset()
	headers[month].WriteString(temp)
	if m.monthInput.Value() != "" || m.focus == month {
		headerStyles[month] = activeHeaderStyle
	} else {
		headerStyles[month] = inactiveHeaderStyle
	}

	headers[dayOfWeek].WriteString("day")
	headers[dayOfWeek].WriteRune('\n')
	headers[dayOfWeek].WriteString("of week")
	temp = lipgloss.PlaceHorizontal(width, lipgloss.Center, headers[dayOfWeek].String())
	headers[dayOfWeek].Reset()
	headers[dayOfWeek].WriteString(temp)
	if m.dayOfWeekInput.Value() != "" || m.focus == dayOfWeek {
		headerStyles[dayOfWeek] = activeHeaderStyle
	} else {
		headerStyles[dayOfWeek] = inactiveHeaderStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		headerStyles[minute].Render(headers[minute].String()),
		headerStyles[hour].Render(headers[hour].String()),
		headerStyles[dayOfMonth].Render(headers[dayOfMonth].String()),
		headerStyles[month].Render(headers[month].String()),
		headerStyles[dayOfWeek].Render(headers[dayOfWeek].String()),
	)
}
