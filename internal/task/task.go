package task

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	markdownID        = regexp.MustCompile(`(^|\s)\$(\d+)(\s|$)`)
	markdownStatus    = regexp.MustCompile(`(?:^|\s)@(next|waiting|parking|someday|review)(?:\s|$)`)
	markdownTags      = regexp.MustCompile(`(?:^|\s)#([A-Za-z0-9][A-Za-z0-9_-]*)`)
	markdownProjects  = regexp.MustCompile(`(?:^|\s)\+([A-Za-z0-9][A-Za-z0-9_-]*)`)
	markdownTimestamp = regexp.MustCompile(`(?:^|\s)(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2}))$`)
	markdownTask      = regexp.MustCompile(`^(\s*)-\s+\[([ xX])\]\s+(.*?)\s*$`)
	markdownComment   = regexp.MustCompile(`^(\s+)-\s+(.*?)\s*$`)
	markdownHeading   = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*#*\s*$`)
)

type Comment struct {
	Text       string
	Tags       []string
	Projects   []string
	LineNumber int
}

type Task struct {
	ID          string
	Description string
	Status      string
	Done        bool
	UpdatedAt   string
	Section     string
	Sections    []string
	Depth       int
	ParentID    string
	LineNumber  int
	Tags        []string
	Projects    []string
	Comments    []Comment
	Line        string
}

type section struct {
	level int
	title string
}

type taskLevel struct {
	indent int
	index  int
}

func FormatTimestamp(value time.Time) string {
	return value.UTC().Format("2006-01-02T15:04:05.000Z")
}

func ParseLines(lines []string) []Task {
	tasks := make([]Task, 0, len(lines))
	sections := make([]section, 0)
	taskLevels := make([]taskLevel, 0)

	for lineIndex, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if match := markdownTask.FindStringSubmatch(line); len(match) == 4 {
			indent := len(match[1])
			for len(taskLevels) > 0 && taskLevels[len(taskLevels)-1].indent >= indent {
				taskLevels = taskLevels[:len(taskLevels)-1]
			}

			item := parseMarkdownTask(line, match[3], match[2] == "x" || match[2] == "X", lineIndex+1)
			item.Depth = len(taskLevels)
			if len(taskLevels) > 0 {
				item.ParentID = tasks[taskLevels[len(taskLevels)-1].index].ID
			}
			item.Sections = sectionTitles(sections)
			if len(item.Sections) > 0 {
				item.Section = item.Sections[len(item.Sections)-1]
			}
			tasks = append(tasks, item)
			taskLevels = append(taskLevels, taskLevel{indent: indent, index: len(tasks) - 1})
			continue
		}

		if match := markdownComment.FindStringSubmatch(line); len(match) == 3 && len(taskLevels) > 0 {
			indent := len(match[1])
			for i := len(taskLevels) - 1; i >= 0; i-- {
				if taskLevels[i].indent < indent {
					tasks[taskLevels[i].index].Comments = append(
						tasks[taskLevels[i].index].Comments,
						parseMarkdownComment(match[2], lineIndex+1),
					)
					break
				}
			}
			continue
		}

		if heading, level, ok := SectionHeading(line); ok {
			for len(sections) > 0 && sections[len(sections)-1].level >= level {
				sections = sections[:len(sections)-1]
			}
			sections = append(sections, section{level: level, title: heading})
			taskLevels = taskLevels[:0]
		}
	}

	return tasks
}

func SectionHeading(line string) (string, int, bool) {
	if match := markdownHeading.FindStringSubmatch(line); len(match) == 3 {
		return strings.TrimSpace(strings.TrimSuffix(match[2], ":")), len(match[1]), true
	}
	trimmed := strings.TrimSpace(line)
	if strings.HasSuffix(trimmed, ":") && !strings.HasPrefix(trimmed, "- ") {
		return strings.TrimSpace(strings.TrimSuffix(trimmed, ":")), 1, true
	}
	return "", 0, false
}

func parseMarkdownTask(line, content string, done bool, lineNumber int) Task {
	item := Task{
		Done:       done,
		Line:       line,
		LineNumber: lineNumber,
	}
	if match := markdownID.FindStringSubmatch(content); len(match) == 4 {
		item.ID = match[2]
	}
	if matches := markdownStatus.FindAllStringSubmatch(content, -1); len(matches) > 0 {
		item.Status = strings.ToLower(matches[len(matches)-1][1])
	}
	item.Tags = captureValues(markdownTags, content)
	item.Projects = captureValues(markdownProjects, content)

	description := content
	if match := markdownTimestamp.FindStringSubmatchIndex(content); len(match) == 4 {
		if updated, err := time.Parse(time.RFC3339Nano, content[match[2]:match[3]]); err == nil {
			item.UpdatedAt = updated.UTC().Format(time.RFC3339Nano)
			description = content[:match[0]] + content[match[1]:]
		}
	}

	description = markdownID.ReplaceAllString(description, "$1$3")
	description = markdownStatus.ReplaceAllString(description, " ")
	description = markdownTags.ReplaceAllString(description, " ")
	description = markdownProjects.ReplaceAllString(description, " ")
	item.Description = strings.Join(strings.Fields(description), " ")
	return item
}

func SetDone(line string, done bool) (string, error) {
	match := markdownTask.FindStringSubmatchIndex(line)
	if len(match) != 8 {
		return "", fmt.Errorf("line is not a Markdown task")
	}
	marker := byte(' ')
	if done {
		marker = 'x'
	}
	return line[:match[4]] + string(marker) + line[match[5]:], nil
}

func SetStatus(line, status string) (string, error) {
	match := markdownTask.FindStringSubmatchIndex(line)
	if len(match) != 8 {
		return "", fmt.Errorf("line is not a Markdown task")
	}
	content := line[match[6]:match[7]]
	statusMatches := markdownStatus.FindAllStringSubmatchIndex(content, -1)
	if len(statusMatches) > 0 {
		for i := len(statusMatches) - 1; i >= 0; i-- {
			statusMatch := statusMatches[i]
			start := statusMatch[2] - 1
			end := statusMatch[3]
			if status == "" || i != len(statusMatches)-1 {
				if start > 0 && isSpace(content[start-1]) {
					start--
				}
				content = content[:start] + content[end:]
				continue
			}
			content = content[:statusMatch[2]] + status + content[end:]
		}
	} else if status != "" && finalTimestampIndex(content) != nil {
		timestamp := finalTimestampIndex(content)
		content = content[:timestamp[0]] + " @" + status + content[timestamp[0]:]
	} else if status != "" {
		content = strings.TrimRight(content, " \t") + " @" + status
	}
	return line[:match[6]] + content + line[match[7]:], nil
}

func isSpace(value byte) bool {
	return value == ' ' || value == '\t'
}

func SetUpdatedAt(line string, updated time.Time) (string, error) {
	match := markdownTask.FindStringSubmatchIndex(line)
	if len(match) != 8 {
		return "", fmt.Errorf("line is not a Markdown task")
	}
	content := line[match[6]:match[7]]
	if timestamp := finalTimestampIndex(content); timestamp != nil {
		content = content[:timestamp[0]]
	} else {
		content = strings.TrimRight(content, " \t")
	}
	content += " " + FormatTimestamp(updated)
	return line[:match[6]] + content + line[match[7]:], nil
}

func finalTimestampIndex(content string) []int {
	return markdownTimestamp.FindStringIndex(content)
}

func parseMarkdownComment(content string, lineNumber int) Comment {
	comment := Comment{
		Tags:       captureValues(markdownTags, content),
		Projects:   captureValues(markdownProjects, content),
		LineNumber: lineNumber,
	}
	text := markdownTags.ReplaceAllString(content, " ")
	text = markdownProjects.ReplaceAllString(text, " ")
	comment.Text = strings.Join(strings.Fields(text), " ")
	return comment
}

func captureValues(pattern *regexp.Regexp, value string) []string {
	matches := pattern.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return nil
	}
	values := make([]string, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		value := match[1]
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return values
}

func sectionTitles(sections []section) []string {
	if len(sections) == 0 {
		return nil
	}
	titles := make([]string, len(sections))
	for i, item := range sections {
		titles[i] = item.title
	}
	return titles
}
