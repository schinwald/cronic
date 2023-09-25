package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Black = "0"
	Red = "1"
	Green = "2"
	Yellow = "3"
	Blue = "4"
	Magenta = "5"
	Cyan = "6"
	White = "7"
	BrightBlack = "8"
	BrightRed = "9"
	BrightGreen = "10"
	BrightYellow = "11"
	BrightBlue = "12"
	BrightMagenta = "13"
	BrightCyan = "14"
	BrightWhite = "15"
)

var (
	PrimaryColor    = lipgloss.Color(Green)
	SecondaryColor  = lipgloss.Color(BrightBlack)
	ForegroundColor = lipgloss.Color(White)
	DimmedForegroundColor = lipgloss.Color(BrightBlack)
	BorderStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(DimmedForegroundColor)
)

func BlockStyle(width int, height int) string {
	var view strings.Builder

	blockStyle := lipgloss.NewStyle().
		Border(lipgloss.BlockBorder()).
		Border(lipgloss.Border{
			Top: "▅",
			Bottom: "🮄",
			Left: "█",
			Right: "█",
			TopRight: "▅",
			BottomRight: "🮄",
			TopLeft: "▅",
			BottomLeft: "🮄",
		}).
		Width(width).
		Height(height).
		Background(lipgloss.Color(White)).
		BorderForeground(lipgloss.Color(White))

	view.WriteString(blockStyle.Render(""))

	return view.String()
}

func PanelStyle(title string, content string, width int, height int, paddingY int, paddingX int) string {
	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Padding(0, 1).
		Render(title)

	boxBorderStyle := BorderStyle.Copy().
		BorderTop(false).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderForeground(DimmedForegroundColor).
		Width(width - 2).
		Height(height - 1).
		Padding(paddingY, paddingX).
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

	boxTopLeftBorderStyle := lipgloss.NewStyle().Foreground(DimmedForegroundColor).Render(boxTopLeftBorder.String())
	boxTopRightBorderStyle := lipgloss.NewStyle().Foreground(DimmedForegroundColor).Render(boxTopRightBorder.String())

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
