package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	Init() tea.Cmd
	Update(tea.Msg) (Page, tea.Cmd)
	View() string
	Size(width int, height int)
}
