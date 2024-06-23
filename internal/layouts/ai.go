package layouts

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/styles"
	"github.com/schinwald/cronic/internal/utils"
)

type AIModel struct {
	explanationLoader        spinner.Model
	explanationInput         components.InputModel
	expression               string
	previous                 string
	current                  string
	debounce                 int
	timeSinceLastInput       int
	naturalLanguageProcessor *utils.NaturalLanguageProcessor
	width                    int
	height                   int
	err                      error
	isGenerating             bool
}

func MakeAIModel() *AIModel {
	explanationLoader := spinner.New()
	explanationLoader.Spinner = spinner.Dot
	explanationLoader.Style = lipgloss.NewStyle().Foreground(styles.TernaryColor)

	explanationInput := components.MakeInputModel()

	naturalLanguageProcessor := utils.MakeNaturalLanguageProcessor()

	return &AIModel{
		explanationLoader:        explanationLoader,
		explanationInput:         explanationInput,
		current:                  "",
		debounce:                 3,
		naturalLanguageProcessor: naturalLanguageProcessor,
	}
}

func (m *AIModel) Init() tea.Cmd {
	m.naturalLanguageProcessor.OnGenerated(func() {
		m.isGenerating = false
	})

	return m.explanationLoader.Tick
}

func (m *AIModel) Update(msg tea.Msg) (*AIModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.explanationInput, cmd = m.explanationInput.Update(msg)
	cmds = append(cmds, cmd)

	// If there is new input reset time since last input
	// Else continue increment the time since last input
	m.previous = m.current
	m.current = m.explanationInput.Value()
	if m.current != m.previous {
		m.timeSinceLastInput = 0

		if !m.isGenerating {
			m.isGenerating = true
			cmds = append(cmds, m.explanationLoader.Tick)
		}
	} else {
		m.timeSinceLastInput++
	}

	// Make a call to the LLM
	if m.timeSinceLastInput == m.debounce {
		m.cancelGeneration()
		go m.generateCronExpressionFromExplanation()
	}

	if m.isGenerating {
		m.explanationLoader, cmd = m.explanationLoader.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *AIModel) View() string {
	var view strings.Builder

	paddingY, paddingX := 1, 3

	var explanation string

	if m.isGenerating {
		explanation = m.explanationLoader.View() + m.explanationInput.View()
	} else {
		explanation = "> " + m.explanationInput.View()
	}

	view.WriteString(styles.PanelStyle("Schedule", explanation, m.width, m.height, paddingY, paddingX))

	return view.String()
}

func (m *AIModel) Size(width int, height int) {
	m.width = width
	m.height = height
}

func (m *AIModel) Focus() error {
	m.explanationInput.Focus()
	return nil
}

func (m *AIModel) Blur() error {
	m.explanationInput.Blur()
	return nil
}

func (m *AIModel) cancelGeneration() {
	m.naturalLanguageProcessor.Cancel()
}

func (m *AIModel) generateCronExpressionFromExplanation() {
	expression, err := m.naturalLanguageProcessor.TextToCronjobExpression(m.explanationInput.Value())
	if err != nil {
		return
	}
	m.expression = expression
}

func (m *AIModel) CronExpression() string {
	return m.expression
}
