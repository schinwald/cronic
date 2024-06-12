package layouts

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/styles"
)

type NextOccurrenceModel struct {
	date   string
	time   string
	width  int
	height int
	err    error
}

func MakeNextOccurrenceModel() NextOccurrenceModel {
	return NextOccurrenceModel{
		date: "N/A",
		time: "N/A",
	}
}

func (m NextOccurrenceModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m NextOccurrenceModel) Update(msg tea.Msg) (NextOccurrenceModel, tea.Cmd) {
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

	headerStyle := lipgloss.NewStyle()
	bodyStyle := lipgloss.NewStyle().Faint(true)

	content.WriteString(headerStyle.Render("Date: "))
	content.WriteString(bodyStyle.Render(m.date))
	content.WriteRune('\n')

	content.WriteString(headerStyle.Render("Time: "))
	content.WriteString(bodyStyle.Render(m.time))

	paddingY, paddingX := 1, 3
	headerBar := styles.PanelStyle("First Occurrence", content.String(), 24, m.height, paddingY, paddingX)

	progressBar := styles.BlockStyle(m.width-lipgloss.Width(headerBar)-2, m.height)

	joined := lipgloss.JoinHorizontal(lipgloss.Top, headerBar, progressBar)
	view.WriteString(joined)

	return view.String()
}

func (m *NextOccurrenceModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *NextOccurrenceModel) NextOccurrence(time time.Time) {
	if time.IsZero() {
		m.date = "N/A"
		m.time = "N/A"
		return
	}

	nextOccurrence := strings.Split(time.String(), " ")
	m.date = nextOccurrence[0]
	m.time = nextOccurrence[1]
}
