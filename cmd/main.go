package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schinwald/cronic/internal/components"
)

type tabModel struct {
	Tabs []string
	TabContent []*components.InputModel
	activeTab int
}

func (tm tabModel) Init() tea.Cmd {
	return nil
}

func (tm tabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return tm, tea.Quit
		case "right", "l", "n", "tab":
			tm.activeTab = min(tm.activeTab+1, len(tm.Tabs)-1)
			return tm, nil
		case "left", "h", "p", "shift+tab":
			tm.activeTab = max(tm.activeTab-1, 0)
			return tm, nil
		default:
			tm.TabContent[tm.activeTab].Update(msg)
			return tm, nil
		}
	}

	return tm, nil
}

func tabBorders(left, middle, right string) lipgloss.Border {
	border := lipgloss.NormalBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorders("┴", "─", "┴")
	activeTabBorder   = tabBorders("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#f76707", Dark: "#f76707"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 2)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true).Bold(true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (tm tabModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, tab := range tm.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tm.Tabs)-1, i == tm.activeTab

		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := style.GetBorder()

		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	input := tm.TabContent[tm.activeTab]

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(input.View()))
	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {	
	tabs := []string{"List", "Timeline", "Details"}
	tabContent := []*components.InputModel{new(components.InputModel), new(components.InputModel), new(components.InputModel)}
	tm := tabModel{Tabs: tabs, TabContent: tabContent}

	if _, err := tea.NewProgram(tm).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
