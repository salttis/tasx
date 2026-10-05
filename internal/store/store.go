package store

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/salttis/tasx/internal/task"
)

var metadataValue = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func Read(path string) ([]task.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open task list %q: %w", path, err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read task list %q: %w", path, err)
	}
	return task.ParseLines(lines), nil
}

func Add(path, project, description string, tags, projects []string, status string) (string, error) {
	description = strings.TrimSpace(description)
	if description == "" || strings.ContainsAny(description, "\r\n") {
		return "", fmt.Errorf("task description must be a non-empty single line")
	}
	if !metadataValue.MatchString(project) {
		return "", fmt.Errorf("invalid project name %q", project)
	}
	if err := validateStatus(status); err != nil {
		return "", err
	}
	for _, tag := range tags {
		if !metadataValue.MatchString(tag) {
			return "", fmt.Errorf("invalid tag %q", tag)
		}
	}
	for _, project := range projects {
		if !metadataValue.MatchString(project) {
			return "", fmt.Errorf("invalid project tag %q", project)
		}
	}

	var id string
	_, err := mutate(path, func(lines *[]string) (bool, error) {
		existingTasks := task.ParseLines(*lines)
		next, err := nextTaskNumber(existingTasks, project)
		if err != nil {
			return false, err
		}
		id = next

		line := "- [ ] " + description
		parsed := task.ParseLines([]string{line})[0]
		for _, tag := range tags {
			if !contains(parsed.Tags, tag) {
				line += " #" + tag
			}
		}
		if status != "" && parsed.Status == "" {
			line += " @" + status
		}
		for _, project := range projects {
			if !contains(parsed.Projects, project) {
				line += " +" + project
			}
		}
		line += " $" + id
		line, err = task.SetUpdatedAt(line, time.Now())
		if err != nil {
			return false, err
		}
		archiveStart, _ := archiveSection(*lines)
		if archiveStart < 0 {
			*lines = append(*lines, line)
		} else {
			before := (*lines)[:archiveStart]
			after := append([]string(nil), (*lines)[archiveStart:]...)
			if len(before) > 0 && strings.TrimSpace(before[len(before)-1]) != "" {
				before = append(before, "")
			}
			before = append(before, line, "")
			*lines = append(before, after...)
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func SetDone(path, project, id string, done bool) (bool, error) {
	return mutateTask(path, project, id, func(item task.Task) (string, bool, error) {
		return updateDone(item, done)
	})
}

func SetDoneAtLine(path string, expected task.Task, done bool) (bool, error) {
	return mutateTaskAtLine(path, expected, func(item task.Task) (string, bool, error) {
		return updateDone(item, done)
	})
}

func SetStatus(path, project, id, status string) (bool, error) {
	if err := validateStatus(status); err != nil {
		return false, err
	}
	return mutateTask(path, project, id, func(item task.Task) (string, bool, error) {
		return updateStatus(item, status)
	})
}

func SetStatusAtLine(path string, expected task.Task, status string) (bool, error) {
	if err := validateStatus(status); err != nil {
		return false, err
	}
	return mutateTaskAtLine(path, expected, func(item task.Task) (string, bool, error) {
		return updateStatus(item, status)
	})
}

func Archive(path, project, id string) (bool, error) {
	return mutate(path, func(lines *[]string) (bool, error) {
		tasks := task.ParseLines(*lines)
		var selected *task.Task
		byLine := make(map[int]task.Task, len(tasks))
		for _, item := range tasks {
			byLine[item.LineNumber] = item
			if !matchesTaskID(item.ID, id, project) {
				continue
			}
			if selected != nil {
				return false, fmt.Errorf("task ID %q is not unique", id)
			}
			copy := item
			selected = &copy
		}
		if selected == nil {
			return false, fmt.Errorf("task ID %q was not found", id)
		}
		if len(selected.Sections) > 0 && strings.EqualFold(selected.Sections[len(selected.Sections)-1], "Arkisto") {
			return false, nil
		}

		start := selected.LineNumber - 1
		end := len(*lines)
		for lineIndex := start + 1; lineIndex < len(*lines); lineIndex++ {
			if _, _, isHeading := task.SectionHeading((*lines)[lineIndex]); isHeading {
				end = lineIndex
				break
			}
			if following, exists := byLine[lineIndex+1]; exists && following.Depth <= selected.Depth {
				end = lineIndex
				break
			}
		}
		block := append([]string(nil), (*lines)[start:end]...)
		updatedLine, err := task.SetUpdatedAt(block[0], time.Now())
		if err != nil {
			return false, err
		}
		block[0] = updatedLine
		indent := len(selected.Line) - len(strings.TrimLeft(selected.Line, " \t"))
		for i, line := range block {
			if line == "" {
				continue
			}
			leading := len(line) - len(strings.TrimLeft(line, " \t"))
			remove := min(indent, leading)
			block[i] = line[remove:]
		}
		for len(block) > 0 && strings.TrimSpace(block[len(block)-1]) == "" {
			block = block[:len(block)-1]
		}
		*lines = append((*lines)[:start], (*lines)[end:]...)

		archiveStart, archiveEnd := archiveSection(*lines)
		if archiveStart < 0 {
			if len(*lines) > 0 && strings.TrimSpace((*lines)[len(*lines)-1]) != "" {
				*lines = append(*lines, "")
			}
			*lines = append(*lines, "Arkisto:")
			*lines = append(*lines, block...)
			return true, nil
		}

		insertAt := archiveEnd
		before := (*lines)[:insertAt]
		after := append([]string(nil), (*lines)[insertAt:]...)
		if len(before) > 0 && strings.TrimSpace(before[len(before)-1]) != "" {
			before = append(before, "")
		}
		before = append(before, block...)
		if len(after) > 0 && strings.TrimSpace(after[0]) != "" &&
			strings.TrimSpace(before[len(before)-1]) != "" {
			before = append(before, "")
		}
		*lines = append(before, after...)
		return true, nil
	})
}

func mutateTask(path, project, id string, update func(task.Task) (string, bool, error)) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, fmt.Errorf("task ID cannot be empty")
	}
	return mutate(path, func(lines *[]string) (bool, error) {
		var selected *task.Task
		for _, item := range task.ParseLines(*lines) {
			if !matchesTaskID(item.ID, id, project) {
				continue
			}
			if selected != nil {
				return false, fmt.Errorf("task ID %q is not unique", id)
			}
			copy := item
			selected = &copy
		}
		if selected == nil {
			return false, fmt.Errorf("task ID %q was not found", id)
		}
		line, changed, err := update(*selected)
		if err != nil || !changed {
			return false, err
		}
		(*lines)[selected.LineNumber-1] = line
		return true, nil
	})
}

func mutateTaskAtLine(path string, expected task.Task, update func(task.Task) (string, bool, error)) (bool, error) {
	return mutate(path, func(lines *[]string) (bool, error) {
		index := expected.LineNumber - 1
		if index < 0 || index >= len(*lines) || (*lines)[index] != expected.Line {
			return false, fmt.Errorf("selected task changed; no changes were saved")
		}
		items := task.ParseLines([]string{(*lines)[index]})
		if len(items) != 1 {
			return false, fmt.Errorf("selected line is no longer a task")
		}
		line, changed, err := update(items[0])
		if err != nil || !changed {
			return false, err
		}
		(*lines)[index] = line
		return true, nil
	})
}

func updateDone(item task.Task, done bool) (string, bool, error) {
	if item.Done == done {
		return item.Line, false, nil
	}
	line, err := task.SetDone(item.Line, done)
	if err != nil {
		return "", false, err
	}
	line, err = task.SetUpdatedAt(line, time.Now())
	return line, true, err
}

func updateStatus(item task.Task, status string) (string, bool, error) {
	if item.Status == status {
		return item.Line, false, nil
	}
	line, err := task.SetStatus(item.Line, status)
	if err != nil {
		return "", false, err
	}
	line, err = task.SetUpdatedAt(line, time.Now())
	return line, true, err
}

func mutate(path string, edit func(*[]string) (bool, error)) (changed bool, err error) {
	lockedPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("resolve task-list path %q: %w", path, err)
	}
	lockedPath = filepath.Clean(lockedPath)
	if runtime.GOOS == "windows" {
		lockedPath = strings.ToLower(lockedPath)
	}
	hash := sha256.Sum256([]byte(lockedPath))
	lockDir := filepath.Join(os.TempDir(), "tasx-locks")
	if err := os.MkdirAll(lockDir, 0o700); err != nil {
		return false, fmt.Errorf("create task-list lock directory: %w", err)
	}
	lock := flock.New(filepath.Join(lockDir, hex.EncodeToString(hash[:])+".lock"))
	if err := lock.Lock(); err != nil {
		return false, fmt.Errorf("lock task list %q: %w", path, err)
	}
	defer func() {
		if unlockErr := lock.Unlock(); err == nil && unlockErr != nil {
			err = fmt.Errorf("unlock task list %q: %w", path, unlockErr)
		}
	}()

	original, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("read task list %q for update: %w", path, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("inspect task list %q for update: %w", path, err)
	}
	text := string(original)
	bom := ""
	if strings.HasPrefix(text, "\uFEFF") {
		bom = "\uFEFF"
		text = strings.TrimPrefix(text, bom)
	}
	newline := "\n"
	if strings.Contains(text, "\r\n") {
		newline = "\r\n"
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	trailingNewline := strings.HasSuffix(text, "\n")
	if trailingNewline {
		text = strings.TrimSuffix(text, "\n")
	}
	var lines []string
	if text != "" {
		lines = strings.Split(text, "\n")
	}
	changed, err = edit(&lines)
	if err != nil || !changed {
		return changed, err
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("recheck task list %q before update: %w", path, err)
	}
	if string(current) != string(original) {
		return false, fmt.Errorf("task list %q changed during update; no changes were saved", path)
	}
	output := strings.Join(lines, newline)
	if trailingNewline || len(lines) > 0 && len(original) == 0 {
		output += newline
	}
	output = bom + output
	if err := writeAtomic(path, []byte(output), info.Mode().Perm()); err != nil {
		return false, fmt.Errorf("save task list %q: %w", path, err)
	}
	return true, nil
}

func nextTaskNumber(tasks []task.Task, _ string) (string, error) {
	maxID := 0
	for _, item := range tasks {
		number, ok := numericTaskID(item.ID)
		if !ok {
			continue
		}
		value, err := strconv.Atoi(number)
		if err != nil {
			return "", fmt.Errorf("parse task number in ID %q: %w", item.ID, err)
		}
		if value > maxID {
			maxID = value
		}
	}
	if maxID == int(^uint(0)>>1) {
		return "", fmt.Errorf("cannot allocate another Tasx task ID")
	}
	return strconv.Itoa(maxID + 1), nil
}

func numericTaskID(id string) (string, bool) {
	if id == "" {
		return "", false
	}
	for _, char := range id {
		if char < '0' || char > '9' {
			return "", false
		}
	}
	return id, true
}

func matchesTaskID(existingID, requestedID, _ string) bool {
	requestedNumber, numeric := numericTaskID(requestedID)
	if !numeric {
		return false
	}
	existingNumber, existingNumeric := numericTaskID(existingID)
	if existingNumeric {
		left, leftErr := strconv.Atoi(existingNumber)
		right, rightErr := strconv.Atoi(requestedNumber)
		return leftErr == nil && rightErr == nil && left == right
	}
	return false
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if strings.EqualFold(value, wanted) {
			return true
		}
	}
	return false
}

func archiveSection(lines []string) (start, end int) {
	start = -1
	level := 0
	for i, line := range lines {
		title, headingLevel, ok := task.SectionHeading(line)
		if ok && strings.EqualFold(title, "Arkisto") {
			start, level = i, headingLevel
			break
		}
	}
	if start < 0 {
		return -1, -1
	}
	end = len(lines)
	for i := start + 1; i < len(lines); i++ {
		if _, headingLevel, ok := task.SectionHeading(lines[i]); ok && headingLevel <= level {
			end = i
			break
		}
	}
	return start, end
}

func validateStatus(status string) error {
	if status == "" || status == "next" || status == "waiting" || status == "parking" || status == "someday" || status == "review" {
		return nil
	}
	return fmt.Errorf("status must be next, waiting, parking, someday, review, or none")
}

func writeAtomic(path string, data []byte, mode os.FileMode) error {
	temp, err := os.CreateTemp(filepath.Dir(path), ".tasx-write-*.tmp")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}
