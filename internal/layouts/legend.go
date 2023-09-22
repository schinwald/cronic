package layouts

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/styles"
)

const (
	minute = iota
	hour
	dayOfMonth
	month
	dayOfWeek
)

type LegendModel struct {
	state int
	err   error
}

func MakeLegendModel() LegendModel {
	return LegendModel{
		state: minute,
	}
}

func (m LegendModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LegendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m LegendModel) View() string {
	var view strings.Builder
	var content strings.Builder

	switch m.state {
	case minute:
		headerStyle := lipgloss.NewStyle().Underline(true).Render("minute")

		var body strings.Builder

		body.WriteString("*  Every minute")
		body.WriteRune('\n')
		body.WriteString(",  At minute x, y, (and ...)")
		body.WriteRune('\n')
		body.WriteString("-  Between x and y")
		body.WriteRune('\n')
		body.WriteString("/  At minute x and every y minutes after")

		bodyStyle := lipgloss.NewStyle().Faint(true).Render(body.String())

		content.WriteString(headerStyle)
		content.WriteRune('\n')
		content.WriteString(bodyStyle)
	case hour:
		content.WriteString("hour")
	case dayOfMonth:
		content.WriteString("day of the month")
	case month:
		content.WriteString("month")
	case dayOfWeek:
		content.WriteString("day of the week")
	}

	view.WriteString(styles.PanelStyle("Legend", content.String(), 100, 15))

	return view.String()
}
