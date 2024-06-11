package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/schinwald/cronic/internal/pages"
	"golang.org/x/term"
)

func main() {
	var err error

	// Open cron file with resource clean up
	file, err := os.OpenFile("/etc/crontab", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, goodCronJobs, badCronJobs := scanCronJobs(file)
	fixedCronJobs := fixBadCronJobs(file, content, badCronJobs)

	var mergedCronJobs []cronJob

	for i := range goodCronJobs {
		mergedCronJobs = append(mergedCronJobs, goodCronJobs[i])
	}

	for i := range fixedCronJobs {
		mergedCronJobs = append(mergedCronJobs, fixedCronJobs[i])
	}

	// Create the audit log if it doesn't exist
	_, err = os.OpenFile("audit.log", os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func fixBadCronJobs(file *os.File, content []string, badCronJobs map[int]cronJob) map[int]cronJob {
	var err error

	// Prompt the user to reconcile bad cron jobs
	for i := range badCronJobs {
		badCronJob := badCronJobs[i]

		badCronJob.tag = fmt.Sprintf("CRONIC-%s", uuid.New().String())
		badCronJob.name = "name"
		badCronJob.description = "description"

		// Update content in memory
		content[i] = fmt.Sprintf("%s %s %s\t\t# %s - Name: %s, Description: %s",
			badCronJob.expression,
			badCronJob.name,
			badCronJob.command,
			badCronJob.tag,
			badCronJob.name,
			badCronJob.description,
		)

		// Seek to beginning of file
		_, err = file.Seek(0, 0)
		if err != nil {
			log.Fatal("unable to reset file for rewrite")
		}

		// Save contents periodically to disk
		writer := bufio.NewWriter(file)
		bytes := []byte(strings.Join(content, "\n"))

		_, err = writer.Write(bytes)
		if err != nil {
			log.Fatal("unable to write to file and update cronjob")
		}

		err = writer.Flush()
		if err != nil {
			log.Fatal("unable to flush buffer and udpate cronjob")
		}
	}

	return badCronJobs
}

func scanCronJobs(file *os.File) ([]string, map[int]cronJob, map[int]cronJob) {
	var content []string
	goodCronJobs := make(map[int]cronJob)
	badCronJobs := make(map[int]cronJob)

	unmarkedJobRegularExpression := regexp.MustCompile(`([^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+)\s+([^\s]+)\s+(.+)`)
	markedJobRegularExpression := regexp.MustCompile(`([^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+)\s+([^\s]+)\s+(.+)\s+# (CRONIC-[^\s]+)\s+-\s+Name:\s([^\s]+),\s+Description:\s([^\s]+)`)
	whiteSpaceRegularExpression := regexp.MustCompile(`\s+`)

	// Read contents of cron file
	reader := bufio.NewReader(file)
	line, isPrefix, err := reader.ReadLine()
	lineNumber := 0

	for err == nil {
		if isPrefix {
			log.Fatal("line is too long to parse cronjobs")
		}

		func() {
			markedJob := markedJobRegularExpression.FindStringSubmatch(string(line))

			if len(markedJob) != 0 {
				func() {
					expression := whiteSpaceRegularExpression.ReplaceAllString(markedJob[1], " ")
					user := markedJob[2]
					command := markedJob[3]
					tag := markedJob[4]
					name := markedJob[5]
					description := markedJob[6]

					// Check if cron expression is valid
					time, err := cronexpr.Parse(expression)
					if err != nil {
						return
					}

					// Add job to jobs map so that we can keep track of it
					goodCronJobs[lineNumber] = cronJob{
						expression:  expression,
						time:        *time,
						user:        user,
						command:     command,
						tag:         tag,
						name:        name,
						description: description,
					}
				}()

				return
			}

			unmarkedJob := unmarkedJobRegularExpression.FindStringSubmatch(string(line))

			// Check if a job was found
			if len(unmarkedJob) != 0 {
				func() {
					expression := whiteSpaceRegularExpression.ReplaceAllString(unmarkedJob[1], " ")
					user := unmarkedJob[2]
					command := unmarkedJob[3]

					// Check if cron expression is valid
					time, err := cronexpr.Parse(expression)
					if err != nil {
						return
					}

					// Add job to jobs map so that we can keep track of it
					badCronJobs[lineNumber] = cronJob{
						expression: expression,
						time:       *time,
						user:       user,
						command:    command,
					}
				}()

				return
			}
		}()

		content = append(content, string(line))
		line, isPrefix, err = reader.ReadLine()
		lineNumber++
	}

	return content, goodCronJobs, badCronJobs
}

const (
	add = iota
	list
)

type cronJob struct {
	expression  string
	time        cronexpr.Expression
	user        string
	command     string
	tag         string
	name        string
	description string
}

type windowModel struct {
	state  int
	pages  []pages.Page
	width  int
	height int
	err    error
}

func initialModel() windowModel {
	pages := []pages.Page{
		pages.MakeAddModel(),
	}

	return windowModel{
		state: add,
		pages: pages,
	}
}

func (m windowModel) Init() tea.Cmd {
	return nil
}

func (m windowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.width, m.height, m.err = term.GetSize(0)

	if m.err != nil {
		return m, nil
	}

	// Handle global update events such as closing the program
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Batch(tea.ClearScreen, tea.Quit)
		}
	case error:
		m.err = msg
		return m, nil
	}

	switch m.state {
	// Handle all actions from initial state
	case add:
		m.pages[add], cmd = m.pages[add].Update(msg)
		cmds = append(cmds, cmd)
	// Handle all actions from list state
	case list:
		m.pages[list], cmd = m.pages[list].Update(msg)
		cmds = append(cmds, cmd)
		break
	}

	return m, tea.Batch(cmds...)
}

func (m windowModel) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)

	m.pages[m.state].Size(
		m.width-lipgloss.Width(style.Render("")),
		m.height-lipgloss.Height(style.Render("")),
	)

	return style.Render(m.pages[m.state].View())
}
