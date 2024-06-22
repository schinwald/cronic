package layouts

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/schinwald/cronic/internal/components"
	"github.com/schinwald/cronic/internal/styles"
	"github.com/schinwald/cronic/internal/utils"
)

type AIModel struct {
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
}

func MakeAIModel() *AIModel {
	explanationInput := components.MakeInputModel()

	naturalLanguageProcessor := utils.MakeNaturalLanguageProcessor()

	return &AIModel{
		explanationInput:         explanationInput,
		current:                  "",
		debounce:                 4,
		naturalLanguageProcessor: naturalLanguageProcessor,
	}
}

func (m *AIModel) Init() tea.Cmd {
	return nil
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
	} else {
		m.timeSinceLastInput++
	}

	// Make a call to the LLM
	if m.timeSinceLastInput == m.debounce {
		m.cancelGeneration()
		go m.generateCronExpressionFromExplanation()
	}

	return m, tea.Batch(cmds...)
}

func (m *AIModel) View() string {
	var view strings.Builder

	paddingY, paddingX := 1, 3

	explanation := "> " + m.explanationInput.View()
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
	m.expression = m.naturalLanguageProcessor.TextToCronjobExpression(m.explanationInput.Value())
}

func (m *AIModel) CronExpression() string {
	return m.expression
}
