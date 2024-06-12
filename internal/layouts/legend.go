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
	focus  int
	width  int
	height int
	err    error
}

func MakeLegendModel() *LegendModel {
	return &LegendModel{
		focus: minute,
	}
}

func (m LegendModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LegendModel) Update(msg tea.Msg) (*LegendModel, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		break
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m LegendModel) View() string {
	var view strings.Builder
	var content strings.Builder

	switch m.focus {
	case minute:
		content.WriteString(createMinuteLegend())
	case hour:
		content.WriteString(createHourLegend())
	case dayOfMonth:
		content.WriteString(createDayOfMonthLegend())
	case month:
		content.WriteString(createMonthLegend())
	case dayOfWeek:
		content.WriteString(createDayOfWeekLegend())
	}

	paddingY, paddingX := 1, 3
	view.WriteString(styles.PanelStyle("Legend", content.String(), m.width, m.height, paddingY, paddingX))

	return view.String()
}

func (m *LegendModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *LegendModel) PreviousFocus() bool {
	switch m.focus {
	case minute:
		return false
	case hour:
		m.SetFocus(minute)
		return true
	case dayOfMonth:
		m.SetFocus(hour)
		return true
	case month:
		m.SetFocus(dayOfMonth)
		return true
	case dayOfWeek:
		m.SetFocus(month)
		return true
	}

	return false
}

func (m *LegendModel) NextFocus() bool {
	switch m.focus {
	case minute:
		m.SetFocus(hour)
		return true
	case hour:
		m.SetFocus(dayOfMonth)
		return true
	case dayOfMonth:
		m.SetFocus(month)
		return true
	case month:
		m.SetFocus(dayOfWeek)
		return true
	case dayOfWeek:
		return false
	}

	return false
}

func (m *LegendModel) SetFocus(focus int) error {
	m.focus = focus

	return nil
}

func createMinuteLegend() string {
	var view strings.Builder

	headerStyle := lipgloss.NewStyle().Underline(true)
	bodyStyle := lipgloss.NewStyle().Faint(true)

	var header strings.Builder

	header.WriteString("minute")

	const (
		inputs = iota
		descriptions
	)

	var columns [2]string
	var columnWidths [2]int
	var columnJustification [2]lipgloss.Position

	columnWidths[inputs] = 10
	columnJustification[inputs] = lipgloss.Left
	columns[inputs] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "*"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], ","),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "-"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "/"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "0-59"),
	)

	columnWidths[descriptions] = 20
	columnJustification[descriptions] = lipgloss.Left
	columns[descriptions] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Every minute"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x, y, (and ...)"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Between x and y"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x and every y minutes after"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x"),
	)

	table := lipgloss.JoinHorizontal(lipgloss.Top,
		columns[inputs],
		columns[descriptions],
	)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header.String()),
		"",
		bodyStyle.Render(table),
	))

	return view.String()
}

func createHourLegend() string {
	var view strings.Builder

	headerStyle := lipgloss.NewStyle().Underline(true)
	bodyStyle := lipgloss.NewStyle().Faint(true)

	var header strings.Builder

	header.WriteString("hour")

	const (
		inputs = iota
		descriptions
	)

	var columns [2]string
	var columnWidths [2]int
	var columnJustification [2]lipgloss.Position

	columnWidths[inputs] = 10
	columnJustification[inputs] = lipgloss.Left
	columns[inputs] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "*"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], ","),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "-"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "/"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "0-23"),
	)

	columnWidths[descriptions] = 25
	columnJustification[descriptions] = lipgloss.Left
	columns[descriptions] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Every minute"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x, y, (and ...)"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Between x and y"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x and every y minutes after"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "At minute x"),
	)

	table := lipgloss.JoinHorizontal(lipgloss.Top,
		columns[inputs],
		columns[descriptions],
	)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header.String()),
		"",
		bodyStyle.Render(table),
	))

	return view.String()
}

func createDayOfMonthLegend() string {
	var view strings.Builder

	headerStyle := lipgloss.NewStyle().Underline(true)
	bodyStyle := lipgloss.NewStyle().Faint(true)

	var header strings.Builder

	header.WriteString("day of month")

	const (
		inputs = iota
		descriptions
	)

	var columns [2]string
	var columnWidths [2]int
	var columnJustification [2]lipgloss.Position

	columnWidths[inputs] = 10
	columnJustification[inputs] = lipgloss.Left
	columns[inputs] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "*"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], ","),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "-"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "/"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "1-31"),
	)

	columnWidths[descriptions] = 20
	columnJustification[descriptions] = lipgloss.Left
	columns[descriptions] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Every day of the month"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x, y, (and ...) of the month"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Between day x and y of the month"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x of the month and every y day(s) after"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x of the month"),
	)

	table := lipgloss.JoinHorizontal(lipgloss.Top,
		columns[inputs],
		columns[descriptions],
	)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header.String()),
		"",
		bodyStyle.Render(table),
	))

	return view.String()
}

func createMonthLegend() string {
	var view strings.Builder

	headerStyle := lipgloss.NewStyle().Underline(true)
	bodyStyle := lipgloss.NewStyle().Faint(true)

	var header strings.Builder

	header.WriteString("month")

	const (
		inputs = iota
		descriptions
	)

	var columns [2]string
	var columnWidths [2]int
	var columnJustification [2]lipgloss.Position

	columnWidths[inputs] = 10
	columnJustification[inputs] = lipgloss.Left
	columns[inputs] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "*"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], ","),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "-"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "/"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "1-12"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "JAN-DEC"),
	)

	columnWidths[descriptions] = 25
	columnJustification[descriptions] = lipgloss.Left
	columns[descriptions] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Every month"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On month x, y, (and ...)"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Between month x and y"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On month x and every y months after"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On month x"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On month x"),
	)

	table := lipgloss.JoinHorizontal(lipgloss.Top,
		columns[inputs],
		columns[descriptions],
	)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header.String()),
		"",
		bodyStyle.Render(table),
	))

	return view.String()
}

func createDayOfWeekLegend() string {
	var view strings.Builder

	headerStyle := lipgloss.NewStyle().Underline(true)
	bodyStyle := lipgloss.NewStyle().Faint(true)

	var header strings.Builder

	header.WriteString("day of week")

	const (
		inputs = iota
		descriptions
	)

	var columns [2]string
	var columnWidths [2]int
	var columnJustification [2]lipgloss.Position

	columnWidths[inputs] = 10
	columnJustification[inputs] = lipgloss.Left
	columns[inputs] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "*"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], ","),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "-"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "/"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "0-6"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "SUN-SAT"),
	)

	columnWidths[descriptions] = 20
	columnJustification[descriptions] = lipgloss.Left
	columns[descriptions] = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Every day of the week"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x, y, (and ...) of the week"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "Between day x and y of the week"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x of the week and every y day(s) after"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x of the week"),
		lipgloss.PlaceHorizontal(columnWidths[inputs], columnJustification[inputs], "On day x of the week"),
	)

	table := lipgloss.JoinHorizontal(lipgloss.Top,
		columns[inputs],
		columns[descriptions],
	)

	view.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header.String()),
		"",
		bodyStyle.Render(table),
	))

	return view.String()
}
