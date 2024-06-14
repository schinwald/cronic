package utils

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/lnquy/cron"
)

type CronJob struct {
	HumanReadable       string
	humanReadableEngine *cron.ExpressionDescriptor
	Next                time.Time
	expression          string
	schedule            cronexpr.Expression
	user                string
	command             string
	tag                 string
	name                string
	description         string
}

func MakeCronJob() *CronJob {
	humanReadableEngine, err := cron.NewDescriptor()
	if err != nil {
		log.Fatal("cannot add human readable engine to cronjob")
	}

	return &CronJob{
		tag:                 fmt.Sprintf("CRONIC-%s", uuid.NewString()),
		humanReadableEngine: humanReadableEngine,
	}
}

func (c *CronJob) Expression(value string) error {
	schedule, err := cronexpr.Parse(value)
	if err != nil {
		return err
	}

	humanReadable, err := c.humanReadableEngine.ToDescription(value, cron.Locale_en)
	if err != nil {
		return err
	}

	c.expression = value
	c.schedule = *schedule
	c.Next = schedule.Next(time.Now())
	c.HumanReadable = humanReadable

	return nil
}

func (c *CronJob) User(value string) {
	c.user = value
}

func (c *CronJob) Command(value string) {
	c.command = value
}

func (c *CronJob) Tag(value string) {
	c.tag = value
}

func (c *CronJob) Name(value string) {
	c.name = value
}

func (c *CronJob) Description(value string) {
	c.description = value
}

func (c CronJob) String() string {
	return fmt.Sprintf("%s %s %s\t\t# %s - Name: %s, Description: %s",
		c.expression,
		c.user,
		c.command,
		c.tag,
		c.name,
		c.description,
	)
}

func (c CronJob) WriteCronJob() error {
	var err error

	// Open cron file with resource clean up
	file, err := os.OpenFile("/etc/crontab", os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Seek(0, 2)
	writer := bufio.NewWriter(file)

	// Write cronjob to the system wide file
	_, err = writer.WriteString(fmt.Sprintf("\n%s", c.String()))
	if err != nil {
		return err
	}

	// Flush the cronjob
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func unmarkedJobRegularExpression() *regexp.Regexp {
	return regexp.MustCompile(`([^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+)\s+([^\s]+)\s+(.+)`)
}

func markedJobRegularExpression() *regexp.Regexp {
	return regexp.MustCompile(`([^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+\s+[^\s]+)\s+([^\s]+)\s+(.+)\s+# (CRONIC-[^\s]+)\s+-\s+Name:\s([^\s]+),\s+Description:\s([^\s]+)`)
}

func passwdRegularExpression() *regexp.Regexp {
	return regexp.MustCompile(`(.*):(.*):(.*):(.*):(.*):(.*):(.*)`)
}

func LoadCronJobs() []CronJob {
	var mergedCronJobs []CronJob

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

func GetAllUsers() []string {
	var users []string
	queue := make(PriorityQueue, 0)

	passwdRegularExpression := passwdRegularExpression()
	rootRegularExpression := regexp.MustCompile("/root")
	homeRegularExpression := regexp.MustCompile("/home")

	// Open cron file with resource clean up
	file, err := os.OpenFile("/etc/passwd", os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, isPrefix, err := reader.ReadLine()
	lineNumber := 0

	for err == nil {
		if isPrefix {
			log.Fatal("line is too long to parse cronjobs")
		}

		passwdMatch := passwdRegularExpression.FindStringSubmatch(string(line))

		if len(passwdMatch) != 0 {
			user := passwdMatch[1]
			directory := passwdMatch[6]

			priority := 0

			homeDirectory := homeRegularExpression.FindStringSubmatch(string(directory))
			if len(homeDirectory) != 0 {
				priority = 50
			}

			rootDirectory := rootRegularExpression.FindStringSubmatch(string(directory))
			if len(rootDirectory) != 0 {
				priority = 100
			}

			item := &PriorityItem{
				value:    user,
				priority: priority,
			}

			heap.Push(&queue, item)
		}

		line, isPrefix, err = reader.ReadLine()
		lineNumber++
	}

	for queue.Len() != 0 {
		item := heap.Pop(&queue).(*PriorityItem)
		users = append(users, item.value)
	}

	return users
}

func fixBadCronJobs(file *os.File, content []string, badCronJobs map[int]CronJob) map[int]CronJob {
	var err error

	// Prompt the user to reconcile bad cron jobs
	for i := range badCronJobs {
		badCronJob := badCronJobs[i]

		badCronJob.tag = fmt.Sprintf("CRONIC-%s", uuid.New().String())
		badCronJob.name = "name"
		badCronJob.description = "description"

		// Update content in memory
		content[i] = badCronJob.String()

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

func scanCronJobs(file *os.File) ([]string, map[int]CronJob, map[int]CronJob) {
	var content []string
	goodCronJobs := make(map[int]CronJob)
	badCronJobs := make(map[int]CronJob)

	unmarkedJobRegularExpression := unmarkedJobRegularExpression()
	markedJobRegularExpression := markedJobRegularExpression()
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
					goodCronJobs[lineNumber] = CronJob{
						expression:  expression,
						schedule:    *time,
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
					badCronJobs[lineNumber] = CronJob{
						expression: expression,
						schedule:   *time,
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
