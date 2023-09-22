package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor    = lipgloss.Color("#00FF00")
	SecondaryColor  = lipgloss.Color("#FF0000")
	ForegroundColor = lipgloss.Color("#777777")
	BorderStyle     = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ForegroundColor)
)

func PanelStyle(title string, content string, width int, height int) string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Padding(0, 1).
		Render(title)

	boxBorderStyle := BorderStyle.
		BorderTop(false).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true).
		Width(width).
		Height(height).
		Padding(2, 5).
		Render(content)

	outerWidth := lipgloss.Width(boxBorderStyle)

	var boxTopLeftBorder strings.Builder
	boxTopLeftBorder.WriteRune('┌')
	for i := 0; i < 1; i++ {
		boxTopLeftBorder.WriteRune('─')
	}

	var boxTopRightBorder strings.Builder
	for i := 0; i < outerWidth-(boxTopLeftBorder.Len()/3)-lipgloss.Width(titleStyle)-1; i++ {
		boxTopRightBorder.WriteRune('─')
	}
	boxTopRightBorder.WriteRune('┐')

	boxTopLeftBorderStyle := lipgloss.NewStyle().Foreground(ForegroundColor).Render(boxTopLeftBorder.String())
	boxTopRightBorderStyle := lipgloss.NewStyle().Foreground(ForegroundColor).Render(boxTopRightBorder.String())

	var boxTopBorder strings.Builder
	boxTopBorder.WriteString(boxTopLeftBorderStyle)
	boxTopBorder.WriteString(titleStyle)
	boxTopBorder.WriteString(boxTopRightBorderStyle)

	s.WriteString(boxTopBorder.String())
	s.WriteRune('\n')
	s.WriteString(boxBorderStyle)
	s.WriteRune('\n')

	return s.String()
}
