package roadmap

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Document struct {
	Title string
	Goals []Goal
}

type Goal struct {
	ID                 string
	Title              string
	Outcome            string
	AcceptanceCriteria []Criterion
	Phases             []Phase
}

type Criterion struct {
	Text string
	Done bool
}

type Phase struct {
	ID        string
	Title     string
	TaskIDs   []string
	BlockedBy []string
}

func Read(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, fmt.Errorf("read roadmap %q: %w", path, err)
	}
	document, err := Parse(string(data))
	if err != nil {
		return Document{}, fmt.Errorf("parse roadmap %q: %w", path, err)
	}
	return document, nil
}

func Parse(content string) (Document, error) {
	var document Document
	var goal *Goal
	var phase *Phase

	for lineIndex, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "# ") && document.Title == "":
			document.Title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		case strings.HasPrefix(trimmed, "## goal:"):
			id, title, ok := parseHeader(trimmed, "## goal:")
			if !ok {
				return Document{}, fmt.Errorf("line %d: malformed goal heading", lineIndex+1)
			}
			document.Goals = append(document.Goals, Goal{ID: id, Title: title})
			goal = &document.Goals[len(document.Goals)-1]
			phase = nil
		case strings.HasPrefix(trimmed, "### phase:"):
			if goal == nil {
				return Document{}, fmt.Errorf("line %d: phase has no preceding goal", lineIndex+1)
			}
			id, title, ok := parseHeader(trimmed, "### phase:")
			if !ok {
				return Document{}, fmt.Errorf("line %d: malformed phase heading", lineIndex+1)
			}
			goal.Phases = append(goal.Phases, Phase{ID: id, Title: title})
			phase = &goal.Phases[len(goal.Phases)-1]
		case goal != nil && phase == nil && strings.HasPrefix(trimmed, "Outcome:"):
			goal.Outcome = strings.TrimSpace(strings.TrimPrefix(trimmed, "Outcome:"))
		case phase != nil && strings.HasPrefix(trimmed, "- task:"):
			fields := strings.Fields(trimmed)
			for _, field := range fields {
				switch {
				case strings.HasPrefix(field, "task:"):
					id := strings.TrimPrefix(field, "task:")
					if id == "" {
						return Document{}, fmt.Errorf("line %d: task reference has no ID", lineIndex+1)
					}
					phase.TaskIDs = append(phase.TaskIDs, id)
				case strings.HasPrefix(field, "blocked-by:"):
					id := strings.TrimPrefix(field, "blocked-by:")
					if id == "" {
						return Document{}, fmt.Errorf("line %d: blocked-by reference has no ID", lineIndex+1)
					}
					phase.BlockedBy = append(phase.BlockedBy, id)
				}
			}
		case goal != nil && phase == nil && strings.HasPrefix(trimmed, "- ["):
			criterion, ok := parseCriterion(trimmed)
			if !ok {
				return Document{}, fmt.Errorf("line %d: malformed acceptance criterion", lineIndex+1)
			}
			goal.AcceptanceCriteria = append(goal.AcceptanceCriteria, criterion)
		}
	}
	if len(document.Goals) == 0 {
		return Document{}, fmt.Errorf("no goal sections found")
	}
	return document, nil
}

func WriteText(output io.Writer, document Document) error {
	if document.Title != "" {
		if _, err := fmt.Fprintf(output, "%s\n", document.Title); err != nil {
			return err
		}
	}
	for goalIndex, goal := range document.Goals {
		if goalIndex > 0 {
			if _, err := fmt.Fprintln(output); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(output, "Tavoite: %s — %s\n", goal.ID, goal.Title); err != nil {
			return err
		}
		if goal.Outcome != "" {
			if _, err := fmt.Fprintf(output, "  Tulos: %s\n", goal.Outcome); err != nil {
				return err
			}
		}
		for _, criterion := range goal.AcceptanceCriteria {
			marker := "[ ]"
			if criterion.Done {
				marker = "[x]"
			}
			if _, err := fmt.Fprintf(output, "  %s %s\n", marker, criterion.Text); err != nil {
				return err
			}
		}
		for _, phase := range goal.Phases {
			if _, err := fmt.Fprintf(output, "  Vaihe: %s — %s\n", phase.ID, phase.Title); err != nil {
				return err
			}
			if len(phase.TaskIDs) > 0 {
				if _, err := fmt.Fprintf(output, "    Tehtävät: %s\n", strings.Join(phase.TaskIDs, ", ")); err != nil {
					return err
				}
			}
			if len(phase.BlockedBy) > 0 {
				if _, err := fmt.Fprintf(output, "    Riippuu: %s\n", strings.Join(phase.BlockedBy, ", ")); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parseHeader(line, prefix string) (string, string, bool) {
	remainder := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	id, title, found := strings.Cut(remainder, "|")
	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	return id, title, found && id != "" && title != ""
}

func parseCriterion(line string) (Criterion, bool) {
	if len(line) < 6 || line[0] != '-' || line[1] != ' ' || line[2] != '[' || line[4] != ']' || line[5] != ' ' {
		return Criterion{}, false
	}
	done := line[3] == 'x' || line[3] == 'X'
	if line[3] != ' ' && !done {
		return Criterion{}, false
	}
	return Criterion{Text: strings.TrimSpace(line[6:]), Done: done}, true
}
