package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schinwald/cronic/internal/styles"
)

const listHeight = 14

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(styles.DimmedForegroundColor)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingTop(3)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprint(i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type ListModel struct {
	list   list.Model
	choice string
	focus  int
}

const (
	on = iota
	off
)

func MakeListModel(items []string) ListModel {
	listItems := []list.Item{}

	for _, value := range items {
		listItems = append(listItems, item(value))
	}

	defaultWidth := 20
	defaultHeight := len(listItems)

	if defaultHeight > 5 {
		defaultHeight = 5
	}

	list := list.New(listItems, itemDelegate{}, defaultWidth, defaultHeight+2)
	list.SetShowTitle(false)
	list.SetFilteringEnabled(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	return ListModel{
		list:  list,
		focus: on,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	var cmd tea.Cmd

	if m.focus == on {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m ListModel) View() string {
	return m.list.View()
}

func (m ListModel) Value() string {
	return m.list.SelectedItem().FilterValue()
}

func (m *ListModel) Focus() error {
	m.focus = on
	return nil
}

func (m *ListModel) Blur() error {
	m.focus = off
	return nil
}
