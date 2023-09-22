package layouts

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/schinwald/cronic/internal/styles"
)

type NextOccurrenceModel struct {
	state int
	err   error
}

func MakeNextOccurrenceModel() NextOccurrenceModel {
	return NextOccurrenceModel{
		state: minute,
	}
}

func (m NextOccurrenceModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m NextOccurrenceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m NextOccurrenceModel) View() string {
	var view strings.Builder
	var content strings.Builder

	switch m.state {
	case minute:
		content.WriteString("minute")
	case hour:
		content.WriteString("hour")
	case dayOfMonth:
		content.WriteString("day of the month")
	case month:
		content.WriteString("month")
	case dayOfWeek:
		content.WriteString("day of the week")
	}

	view.WriteString(styles.PanelStyle("Next Occurrence", content.String(), 200, 5))

	return view.String()
}
