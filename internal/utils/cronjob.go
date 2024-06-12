package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
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

func LoadCronJobs() []cronJob {
	var mergedCronJobs []cronJob

	// Open cron file with resource clean up
	file, err := os.OpenFile("/etc/crontab", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, goodCronJobs, badCronJobs := scanCronJobs(file)
	fixedCronJobs := fixBadCronJobs(file, content, badCronJobs)

	for i := range goodCronJobs {
		mergedCronJobs = append(mergedCronJobs, goodCronJobs[i])
	}

	for i := range fixedCronJobs {
		mergedCronJobs = append(mergedCronJobs, fixedCronJobs[i])
	}

	return mergedCronJobs
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
