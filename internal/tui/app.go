package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/salttis/tasx/internal/config"
	"github.com/salttis/tasx/internal/roadmap"
	"github.com/salttis/tasx/internal/scope"
	"github.com/salttis/tasx/internal/store"
	"github.com/salttis/tasx/internal/task"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	focusStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("6"))
	paneStyle  = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type model struct {
	scopes      []scope.Scope
	scopeIndex  int
	tasks       []task.Task
	roadmap     roadmap.Document
	roadmapErr  error
	taskIndex   int
	focus       int
	stateFilter string
	search      string
	searching   bool
	adding      bool
	pending     bool
	message     string
	width       int
	height      int
}

type operationResultMsg struct {
	message string
	err     error
}

func Run(cfg config.Config, forceGlobal, forceRepo bool) error {
	if forceGlobal && forceRepo {
		return fmt.Errorf("valitse vain toinen: --global tai --repo")
	}
	scopes, err := scope.List(cfg)
	if err != nil {
		return err
	}
	m := model{scopes: scopes, stateFilter: cfg.State}
	if forceGlobal || (!forceRepo && cfg.DefaultScope == "global") {
		m.scopeIndex = findScope(scopes, cfg.PersonalDir)
	} else {
		current, err := scope.Repository()
		if err != nil {
			if forceRepo || cfg.DefaultScope == "repo" {
				return err
			}
			m.scopeIndex = findScope(scopes, cfg.PersonalDir)
		} else {
			m.scopeIndex = findScope(scopes, current.Root)
			if m.scopeIndex < 0 {
				m.scopes = append(m.scopes, current)
				m.scopeIndex = len(m.scopes) - 1
			}
		}
	}
	if m.scopeIndex < 0 {
		m.scopeIndex = 0
	}
	m.loadSelected()

	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run Tasx TUI: %w", err)
	}
	return nil
}

func findScope(scopes []scope.Scope, root string) int {
	for i, item := range scopes {
		if samePath(item.Root, root) {
			return i
		}
	}
	return -1
}

func samePath(left, right string) bool {
	return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case operationResultMsg:
		m.pending = false
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			m.loadSelected()
			if m.message == "" {
				m.message = msg.message
			}
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		return m.updateKey(msg)
	}
	return m, nil
}

func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if m.pending {
		return m, nil
	}
	if m.searching {
		switch msg.Type {
		case tea.KeyEsc:
			m.searching = false
		case tea.KeyEnter:
			m.searching = false
		case tea.KeyBackspace:
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
			}
		case tea.KeyRunes:
			m.search += string(msg.Runes)
		}
		m.clampTask()
		return m, nil
	}
	if m.adding {
		switch msg.Type {
		case tea.KeyEsc:
			m.adding = false
			m.search = ""
		case tea.KeyEnter:
			description := strings.TrimSpace(m.search)
			m.adding = false
			m.search = ""
			if description == "" {
				m.message = "Tyhjä tehtävä ohitettiin."
				return m, nil
			}
			return m.addTask(description)
		case tea.KeyBackspace:
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
			}
		case tea.KeyRunes:
			m.search += string(msg.Runes)
		}
		return m, nil
	}
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.message = "Tab fokus | j/k selaa | / hae | a lisää | Space valmis | s tila | f suodatin | r päivitä | q sulje"
	case "tab":
		m.focus = (m.focus + 1) % 3
	case "h":
		m.focus = (m.focus + 2) % 3
	case "l":
		m.focus = (m.focus + 1) % 3
	case "j", "down":
		if m.focus == 0 && len(m.scopes) > 0 {
			m.scopeIndex = (m.scopeIndex + 1) % len(m.scopes)
			m.loadSelected()
		} else if m.focus == 1 {
			m.taskIndex = min(m.taskIndex+1, max(0, len(m.visibleTasks())-1))
		}
	case "k", "up":
		if m.focus == 0 && len(m.scopes) > 0 {
			m.scopeIndex = (m.scopeIndex - 1 + len(m.scopes)) % len(m.scopes)
			m.loadSelected()
		} else if m.focus == 1 {
			m.taskIndex = max(0, m.taskIndex-1)
		}
	case "/":
		m.searching = true
		m.search = ""
		m.focus = 1
	case "f":
		switch m.stateFilter {
		case "all":
			m.stateFilter = "open"
		case "open":
			m.stateFilter = "done"
		default:
			m.stateFilter = "all"
		}
		m.clampTask()
	case "r":
		m.loadSelected()
	case "a":
		if m.scopeIndex >= 0 && m.scopeIndex < len(m.scopes) && !m.scopes[m.scopeIndex].Missing {
			m.adding = true
			m.search = ""
		} else {
			m.message = "Puuttuvaa tehtävälistaa ei voi muokata."
		}
	case " ", "d":
		return m.setSelectedDone()
	case "s":
		return m.cycleSelectedStatus()
	}
	return m, nil
}

func (m model) View() string {
	if len(m.scopes) == 0 {
		return "Tasx: yhtaan scopea ei loytynyt.\n\nq sulkee"
	}
	visible := m.visibleTasks()
	leftWidth, middleWidth, rightWidth := paneWidths(m.width)
	bodyHeight := max(6, m.height-4)

	scopeLines := make([]string, 0, len(m.scopes))
	for i, item := range m.scopes {
		marker := "  "
		if i == m.scopeIndex {
			marker = "> "
		}
		label := item.Name
		if item.Missing {
			label += " (puuttuu)"
		}
		scopeLines = append(scopeLines, marker+label)
	}
	taskLines := make([]string, 0, len(visible)+1)
	for i, item := range visible {
		marker := "  "
		if i == m.taskIndex {
			marker = "> "
		}
		status := " "
		if item.Done {
			status = "x"
		}
		line := fmt.Sprintf("%s%s- [%s] %s", marker, strings.Repeat("  ", item.Depth), status, item.Description)
		if item.Status != "" {
			line += " @" + item.Status
		}
		if item.ID != "" {
			line += " $" + item.ID
		}
		taskLines = append(taskLines, line)
	}
	details := []string{
		"Scope: " + m.scopes[m.scopeIndex].Name,
		"Polku:",
		m.scopes[m.scopeIndex].TasksPath,
	}
	if len(visible) > 0 && m.taskIndex < len(visible) && len(visible[m.taskIndex].Sections) > 0 {
		details = append(details, "Osio: "+strings.Join(visible[m.taskIndex].Sections, " / "))
	}
	if len(visible) > 0 && m.taskIndex < len(visible) {
		selected := visible[m.taskIndex]
		details = append(details, "", "ID: "+selected.ID, "Tila: "+selected.Status)
		if len(selected.Tags) > 0 {
			details = append(details, "Tagit: "+strings.Join(selected.Tags, ", "))
		}
		if len(selected.Projects) > 0 {
			details = append(details, "Projektit: "+strings.Join(selected.Projects, ", "))
		}
		if selected.UpdatedAt != "" {
			if updated, err := time.Parse(time.RFC3339Nano, selected.UpdatedAt); err == nil {
				details = append(details, "Päivitetty: "+task.FormatTimestamp(updated))
			}
		}
		for _, comment := range selected.Comments {
			details = append(details, "", "- "+comment.Text)
		}
		details = append(details, "", selected.Description)
	}
	if selected, ok := m.selectedTask(); ok {
		details = append(details, "")
		details = append(details, m.roadmapDetails(selected)...)
	} else {
		details = append(details, "")
		details = append(details, m.roadmapDetails(nil)...)
	}
	if m.message != "" {
		details = append(details, "", m.message)
	}
	if m.scopes[m.scopeIndex].Missing {
		taskLines = []string{"Tehtävälista puuttuu.", "", m.scopes[m.scopeIndex].TasksPath}
	}
	if len(taskLines) == 0 {
		taskLines = []string{"Ei näytettäviä tehtäviä."}
	}
	filter := fmt.Sprintf("Suodatin: %s", m.stateFilter)
	taskLines = append([]string{filter}, taskLines...)

	var body string
	if m.width > 0 && m.width < 90 {
		titles := []string{"SCOPES", "TASKS", "DETAILS"}
		content := [][]string{scopeLines, taskLines, details}
		body = renderPane(titles[m.focus], content[m.focus], max(24, m.width-2), bodyHeight, true)
	} else {
		panes := []string{
			renderPane("SCOPES", scopeLines, leftWidth, bodyHeight, m.focus == 0),
			renderPane("TASKS", taskLines, middleWidth, bodyHeight, m.focus == 1),
			renderPane("DETAILS", details, rightWidth, bodyHeight, m.focus == 2),
		}
		body = lipgloss.JoinHorizontal(lipgloss.Top, panes...)
	}
	footer := "Tab/h/l fokus  j/k selaa  / haku  a lisää  Space valmis  s tila  f suodatin  r päivitä  q sulje"
	if m.searching {
		footer = "Haku: " + m.search + "  (Enter valmis, Esc peruuta)"
	} else if m.adding {
		footer = "Uusi tehtävä: " + m.search + "  (Enter lisää, Esc peruuta)"
	} else if m.pending {
		footer = "Tallennetaan tehtävämuutosta..."
	}
	return titleStyle.Render(" TASX ") + "\n" +
		body + "\n" +
		mutedStyle.Render(footer)
}

func (m *model) loadSelected() {
	m.message = ""
	m.tasks = nil
	m.roadmap = roadmap.Document{}
	m.roadmapErr = nil
	if m.scopeIndex < 0 || m.scopeIndex >= len(m.scopes) {
		return
	}
	selected := m.scopes[m.scopeIndex]
	document, roadmapErr := roadmap.Read(selected.RoadmapPath)
	if roadmapErr != nil {
		m.roadmapErr = roadmapErr
	} else {
		m.roadmap = document
	}
	if selected.Missing {
		m.message = "Repo-scopea ei vaihdettu automaattisesti."
		return
	}
	tasks, err := store.Read(selected.TasksPath)
	if err != nil {
		m.message = err.Error()
		return
	}
	m.tasks = tasks
	m.clampTask()
}

func (m model) visibleTasks() []task.Task {
	items := make([]task.Task, 0, len(m.tasks))
	for _, item := range m.tasks {
		if m.stateFilter == "open" && item.Done || m.stateFilter == "done" && !item.Done {
			continue
		}
		if m.search != "" && !strings.Contains(strings.ToLower(item.Description), strings.ToLower(m.search)) && !strings.Contains(strings.ToLower(item.ID), strings.ToLower(m.search)) {
			continue
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].LineNumber < items[j].LineNumber
	})
	return items
}

func (m model) selectedTask() (*task.Task, bool) {
	items := m.visibleTasks()
	if m.taskIndex < 0 || m.taskIndex >= len(items) {
		return nil, false
	}
	return &items[m.taskIndex], true
}

func (m model) roadmapDetails(selected *task.Task) []string {
	lines := []string{"ROADMAP"}
	if m.roadmapErr != nil {
		if errors.Is(m.roadmapErr, os.ErrNotExist) {
			return append(lines, "Roadmapia ei ole.")
		}
		return append(lines, "Roadmapin luku epäonnistui: "+m.roadmapErr.Error())
	}
	for _, goal := range m.roadmap.Goals {
		lines = append(lines, goal.Title)
		if goal.Outcome != "" {
			lines = append(lines, goal.Outcome)
		}
		for _, phase := range goal.Phases {
			marker := "  "
			if selected != nil && phaseContainsTask(phase, selected.ID) {
				marker = "> "
			}
			lines = append(lines, marker+phase.Title)
			if selected != nil && marker == "> " && len(phase.TaskIDs) > 0 {
				lines = append(lines, "Tehtävä: "+strings.Join(phase.TaskIDs, ", "))
			}
		}
	}
	if len(m.roadmap.Goals) == 0 {
		lines = append(lines, "Roadmap ei sisällä tavoitteita.")
	}
	return lines
}

func phaseContainsTask(phase roadmap.Phase, taskID string) bool {
	taskNumber, ok := numericTaskID(taskID)
	if !ok {
		return false
	}
	for _, reference := range phase.TaskIDs {
		referenceNumber, ok := numericTaskID(reference)
		if ok && referenceNumber == taskNumber {
			return true
		}
	}
	return false
}

func numericTaskID(id string) (uint64, bool) {
	if separator := strings.LastIndexByte(id, '-'); separator >= 0 {
		id = id[separator+1:]
	}
	number, err := strconv.ParseUint(id, 10, 64)
	return number, err == nil
}

func (m model) addTask(description string) (tea.Model, tea.Cmd) {
	selected, ok := m.selectedScope()
	if !ok {
		m.message = "Valittua tehtävälistaa ei voi muokata."
		return m, nil
	}
	m.pending = true
	project := scope.ProjectKey(selected.Name)
	return m, func() tea.Msg {
		id, err := store.Add(selected.TasksPath, project, description, nil, nil, "")
		return operationResultMsg{message: "Lisättiin tehtävä $" + id + ".", err: err}
	}
}

func (m model) setSelectedDone() (tea.Model, tea.Cmd) {
	selectedScope, scopeOK := m.selectedScope()
	selectedTask, taskOK := m.selectedTask()
	if !scopeOK || !taskOK {
		m.message = "Valitse muokattava tehtävä."
		return m, nil
	}
	done := !selectedTask.Done
	m.pending = true
	project := scope.ProjectKey(selectedScope.Name)
	action := "merkittiin valmiiksi."
	if !done {
		action = "avattiin uudelleen."
	}
	taskLabel := selectedTask.ID
	if taskLabel == "" {
		taskLabel = fmt.Sprintf("rivin %d tehtävä", selectedTask.LineNumber)
	} else {
		taskLabel = "$" + taskLabel
	}
	return m, func() tea.Msg {
		var err error
		if selectedTask.ID == "" {
			_, err = store.SetDoneAtLine(selectedScope.TasksPath, *selectedTask, done)
		} else {
			_, err = store.SetDone(selectedScope.TasksPath, project, selectedTask.ID, done)
		}
		return operationResultMsg{message: "Tehtävä " + taskLabel + " " + action, err: err}
	}
}

func (m model) cycleSelectedStatus() (tea.Model, tea.Cmd) {
	selectedScope, scopeOK := m.selectedScope()
	selectedTask, taskOK := m.selectedTask()
	if !scopeOK || !taskOK {
		m.message = "Valitse muokattava tehtävä."
		return m, nil
	}
	statuses := []string{"", "next", "waiting", "parking", "someday", "review"}
	nextStatus := statuses[0]
	for index, status := range statuses {
		if status == selectedTask.Status {
			nextStatus = statuses[(index+1)%len(statuses)]
			break
		}
	}
	m.pending = true
	project := scope.ProjectKey(selectedScope.Name)
	return m, func() tea.Msg {
		var err error
		if selectedTask.ID == "" {
			_, err = store.SetStatusAtLine(selectedScope.TasksPath, *selectedTask, nextStatus)
		} else {
			_, err = store.SetStatus(selectedScope.TasksPath, project, selectedTask.ID, nextStatus)
		}
		taskLabel := selectedTask.ID
		if taskLabel == "" {
			taskLabel = fmt.Sprintf("rivin %d tehtävän", selectedTask.LineNumber)
		} else {
			taskLabel = "$" + taskLabel
		}
		message := "Tehtävän " + taskLabel + " tila poistettiin."
		if nextStatus != "" {
			message = "Tehtävän " + taskLabel + " tilaksi asetettiin @" + nextStatus + "."
		}
		return operationResultMsg{message: message, err: err}
	}
}

func (m model) selectedScope() (scope.Scope, bool) {
	if m.scopeIndex < 0 || m.scopeIndex >= len(m.scopes) || m.scopes[m.scopeIndex].Missing {
		return scope.Scope{}, false
	}
	return m.scopes[m.scopeIndex], true
}

func (m *model) clampTask() {
	count := len(m.visibleTasks())
	m.taskIndex = min(m.taskIndex, max(0, count-1))
}

func paneWidths(total int) (int, int, int) {
	if total < 90 {
		total = 90
	}
	left := max(22, total/5)
	middle := max(34, total*2/5)
	right := max(28, total-left-middle-4)
	return left, middle, right
}

func renderPane(title string, lines []string, width, height int, focused bool) string {
	content := titleStyle.Render(title) + "\n" + strings.Join(lines, "\n")
	style := paneStyle
	if focused {
		style = focusStyle
	}
	return style.Width(width).Height(height).Render(content)
}
