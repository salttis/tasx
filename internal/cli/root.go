package cli

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/salttis/tasx/internal/config"
	"github.com/salttis/tasx/internal/roadmap"
	"github.com/salttis/tasx/internal/scope"
	"github.com/salttis/tasx/internal/store"
	"github.com/salttis/tasx/internal/task"
	"github.com/salttis/tasx/internal/tui"
	"github.com/spf13/cobra"
)

var version = "dev"

func Execute(txAlias bool) error {
	var cfg config.Config
	var configPath string
	var forceGlobal, forceRepo bool
	var stateFilter string

	root := &cobra.Command{
		Use:          "tasx",
		Short:        "Local-first tehtävienhallinta",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if txAlias {
				return runTUI(cfg, forceGlobal, forceRepo)
			}
			project := ""
			if len(args) > 0 {
				project = args[0]
			}
			return runList(cmd, cfg, forceGlobal, forceRepo, stateFilter, project)
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var err error
			cfg, configPath, err = config.Load()
			return err
		},
	}
	root.PersistentFlags().BoolVar(&forceGlobal, "global", false, "käytä käyttäjätason tehtävälistaa")
	root.PersistentFlags().BoolVar(&forceRepo, "repo", false, "vaadi nykyisen Git-repon tehtävälista")
	root.PersistentFlags().StringVar(&stateFilter, "state", "", "suodata: all, open tai done")

	list := &cobra.Command{
		Use:   "list",
		Short: "Näytä tehtävät",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFilter := ""
			if len(args) > 0 {
				projectFilter = args[0]
			}
			return runList(cmd, cfg, forceGlobal, forceRepo, stateFilter, projectFilter)
		},
	}
	root.AddCommand(list)

	var addStatus, targetProject string
	var addTags []string
	addCommand := &cobra.Command{
		Use:   "add <description>",
		Short: "Lisää tehtävä valittuun tehtävälistaan",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			status := strings.TrimPrefix(addStatus, "@")
			if status == "none" {
				status = ""
			}
			return runAdd(cmd, cfg, forceGlobal, forceRepo, strings.Join(args, " "), addTags, status, targetProject)
		},
	}
	addCommand.Flags().StringSliceVar(&addTags, "tag", nil, "tehtävän tagi (toistettava)")
	addCommand.Flags().StringVar(&targetProject, "projekti", "", "kohdeprojekti, kun tehtävä lisätään toisesta scopesta")
	addCommand.Flags().StringVar(&targetProject, "project", "", "alias lipulle --projekti")
	addCommand.Flags().StringVar(&addStatus, "status", "", "työnkulkutila: next, waiting, parking, someday tai review")
	root.AddCommand(addCommand)

	doneCommand := &cobra.Command{
		Use:   "done <task-id>",
		Short: "Merkitse tehtävä valmiiksi",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetDone(cmd, cfg, forceGlobal, forceRepo, args[0], true)
		},
	}
	root.AddCommand(doneCommand)

	reopenCommand := &cobra.Command{
		Use:   "reopen <task-id>",
		Short: "Avaa valmis tehtävä uudelleen",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetDone(cmd, cfg, forceGlobal, forceRepo, args[0], false)
		},
	}
	root.AddCommand(reopenCommand)

	statusCommand := &cobra.Command{
		Use:   "status <task-id> <status>",
		Short: "Aseta tai poista tehtävän työnkulkutila",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			status := strings.TrimPrefix(args[1], "@")
			if status == "none" {
				status = ""
			}
			return runSetStatus(cmd, cfg, forceGlobal, forceRepo, args[0], status)
		},
	}
	root.AddCommand(statusCommand)

	archiveCommand := &cobra.Command{
		Use:   "archive <task-id>",
		Short: "Siirrä tehtävä Arkisto-osioon käyttäjän pyynnöstä",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchive(cmd, cfg, forceGlobal, forceRepo, args[0])
		},
	}
	root.AddCommand(archiveCommand)

	tuiCommand := &cobra.Command{
		Use:   "tui",
		Short: "Avaa näppäimistökäyttöinen käyttöliittymä",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runTUI(cfg, forceGlobal, forceRepo)
		},
	}
	root.AddCommand(tuiCommand)

	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Näytä käytössä olevat käyttäjäasetukset",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Henkilökohtaisen datan hakemisto: %s\nOletusscope: %s\nTilan oletussuodatin: %s\nAsetustiedosto: %s\n",
				cfg.PersonalDir, cfg.DefaultScope, cfg.State, configPath)
			return err
		},
	}
	root.AddCommand(configCommand)

	roadmapCommand := &cobra.Command{
		Use:   "roadmap [project]",
		Short: "Näytä valitun scopen roadmap",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := ""
			if len(args) > 0 {
				projectName = args[0]
			}
			return runRoadmap(cmd, cfg, forceGlobal, forceRepo, projectName)
		},
	}
	root.AddCommand(roadmapCommand)
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Näytä ohjelmaversio",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		},
	})
	return root.Execute()
}

func runList(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, stateFilter, projectName string) error {
	if forceGlobal && forceRepo {
		return fmt.Errorf("valitse vain toinen: --global tai --repo")
	}
	if projectName != "" && (forceGlobal || forceRepo) {
		return fmt.Errorf("projektin nimellä listattaessa älä käytä --global- tai --repo-valitsinta")
	}
	if stateFilter == "" {
		stateFilter = cfg.State
	}
	if stateFilter != "all" && stateFilter != "open" && stateFilter != "done" {
		return fmt.Errorf("--state must be all, open, or done")
	}

	selected, err := commandScope(cfg, forceGlobal, forceRepo, projectName)
	if err != nil {
		return err
	}
	if selected.Missing {
		return fmt.Errorf("task list does not exist: %s", selected.TasksPath)
	}
	tasks, err := store.Read(selected.TasksPath)
	if err != nil {
		return err
	}
	tasks = filteredAndSorted(tasks, stateFilter)
	if len(tasks) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Ei suodatukseen sopivia tehtäviä.")
		return nil
	}
	printTasks(cmd.OutOrStdout(), tasks, selected.Name)
	return nil
}

func runRoadmap(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, projectName string) error {
	if forceGlobal && forceRepo {
		return fmt.Errorf("valitse vain toinen: --global tai --repo")
	}
	var selected scope.Scope
	var err error
	switch {
	case projectName != "":
		if forceGlobal || forceRepo {
			return fmt.Errorf("projektin nimellä roadmapia haettaessa älä käytä --global- tai --repo-valitsinta")
		}
		selected, err = scope.ResolveProject(cfg, projectName)
	case forceGlobal:
		selected, err = scope.Current(cfg, true, false)
	case forceRepo:
		selected, err = scope.Repository()
	default:
		selected, err = scope.Context(cfg)
	}
	if err != nil {
		return err
	}
	document, err := roadmap.Read(selected.RoadmapPath)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Scope: %s\n", selected.Name); err != nil {
		return err
	}
	return roadmap.WriteText(cmd.OutOrStdout(), document)
}

func runTUI(cfg config.Config, forceGlobal, forceRepo bool) error {
	if forceGlobal && forceRepo {
		return fmt.Errorf("valitse vain toinen: --global tai --repo")
	}
	return tui.Run(cfg, forceGlobal, forceRepo)
}

func runAdd(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, description string, tags []string, status, targetProject string) error {
	selected, err := commandScope(cfg, forceGlobal, forceRepo, targetProject)
	if err != nil {
		return err
	}
	idPrefix := scope.ProjectKey(selected.Name)
	if idPrefix == "" {
		return fmt.Errorf("could not derive project name for scope %q", selected.Name)
	}
	id, err := store.Add(selected.TasksPath, idPrefix, description, tags, nil, status)
	if err != nil {
		return err
	}
	if targetProject != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Lisättiin tehtävä %s-%s (%s)\n", scope.ProjectKey(selected.Name), id, selected.Name)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Lisättiin tehtävä $%s (%s)\n", id, selected.Name)
	}
	return nil
}

func runSetDone(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, id string, done bool) error {
	selected, localID, err := taskScope(cfg, forceGlobal, forceRepo, id)
	if err != nil {
		return err
	}
	changed, err := store.SetDone(selected.TasksPath, scope.ProjectKey(selected.Name), localID, done)
	if err != nil {
		return err
	}
	if !changed {
		if done {
			fmt.Fprintf(cmd.OutOrStdout(), "Tehtävä %s oli jo valmis.\n", id)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Tehtävä %s oli jo avoin.\n", id)
		}
		return nil
	}
	action := "avattiin uudelleen"
	if done {
		action = "merkittiin valmiiksi"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Tehtävä %s %s.\n", id, action)
	return nil
}

func runSetStatus(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, id, status string) error {
	selected, localID, err := taskScope(cfg, forceGlobal, forceRepo, id)
	if err != nil {
		return err
	}
	changed, err := store.SetStatus(selected.TasksPath, scope.ProjectKey(selected.Name), localID, status)
	if err != nil {
		return err
	}
	if !changed {
		fmt.Fprintf(cmd.OutOrStdout(), "Tehtävän %s tila oli jo ennallaan.\n", id)
		return nil
	}
	if status == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Tehtävän %s työnkulkutila poistettiin.\n", id)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Tehtävän %s tilaksi asetettiin @%s.\n", id, status)
	return nil
}

func runArchive(cmd *cobra.Command, cfg config.Config, forceGlobal, forceRepo bool, id string) error {
	selected, localID, err := taskScope(cfg, forceGlobal, forceRepo, id)
	if err != nil {
		return err
	}
	changed, err := store.Archive(selected.TasksPath, scope.ProjectKey(selected.Name), localID)
	if err != nil {
		return err
	}
	if !changed {
		fmt.Fprintf(cmd.OutOrStdout(), "Tehtävä %s on jo Arkisto-osiossa.\n", id)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Tehtävä %s siirrettiin Arkisto-osioon.\n", id)
	return nil
}

func commandScope(cfg config.Config, forceGlobal, forceRepo bool, projectName string) (scope.Scope, error) {
	if forceGlobal && forceRepo {
		return scope.Scope{}, fmt.Errorf("valitse vain toinen: --global tai --repo")
	}
	var selected scope.Scope
	var err error
	switch {
	case projectName != "":
		if forceGlobal || forceRepo {
			return scope.Scope{}, fmt.Errorf("--projekti-valitsinta ei voi yhdistää --global- tai --repo-valitsimeen")
		}
		selected, err = scope.ResolveProject(cfg, projectName)
	case forceGlobal, forceRepo:
		selected, err = scope.Current(cfg, forceGlobal, forceRepo)
	default:
		selected, err = scope.Context(cfg)
	}
	if err != nil {
		return scope.Scope{}, err
	}
	if selected.Missing {
		return scope.Scope{}, fmt.Errorf("task list does not exist: %s", selected.TasksPath)
	}
	return selected, nil
}

func taskScope(cfg config.Config, forceGlobal, forceRepo bool, id string) (scope.Scope, string, error) {
	if prefix, number, qualified := splitQualifiedID(id); qualified {
		if !forceGlobal && !forceRepo {
			current, err := scope.Context(cfg)
			if err != nil {
				return scope.Scope{}, "", err
			}
			if scope.ProjectKey(current.Name) == scope.ProjectKey(prefix) {
				return current, number, nil
			}
		}
		selected, err := commandScope(cfg, forceGlobal, forceRepo, prefix)
		if err == nil {
			return selected, number, nil
		}
		current, currentErr := commandScope(cfg, forceGlobal, forceRepo, "")
		if currentErr == nil && scope.ProjectKey(current.Name) == scope.ProjectKey(prefix) {
			return current, number, nil
		}
		return scope.Scope{}, "", err
	}
	selected, err := commandScope(cfg, forceGlobal, forceRepo, "")
	if err != nil {
		return scope.Scope{}, "", err
	}
	return selected, id, nil
}

func splitQualifiedID(id string) (string, string, bool) {
	index := strings.LastIndexByte(id, '-')
	if index <= 0 || index == len(id)-1 {
		return "", "", false
	}
	number := id[index+1:]
	if _, err := strconv.ParseUint(number, 10, 64); err != nil {
		return "", "", false
	}
	return id[:index], number, true
}

func localTaskNumber(id, project string) (string, bool) {
	if _, err := strconv.ParseUint(id, 10, 64); err == nil {
		return id, true
	}
	index := strings.LastIndexByte(id, '-')
	if index <= 0 || index == len(id)-1 || scope.ProjectKey(id[:index]) != scope.ProjectKey(project) {
		return "", false
	}
	number := id[index+1:]
	_, err := strconv.ParseUint(number, 10, 64)
	return number, err == nil
}

func filteredAndSorted(tasks []task.Task, stateFilter string) []task.Task {
	filtered := make([]task.Task, 0, len(tasks))
	for _, item := range tasks {
		if stateFilter == "open" && item.Done || stateFilter == "done" && !item.Done {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].LineNumber < filtered[j].LineNumber
	})
	return filtered
}

func printTasks(output io.Writer, tasks []task.Task, project string) {
	type taskEvent struct {
		line    int
		task    *task.Task
		comment *task.Comment
		depth   int
	}
	events := make([]taskEvent, 0, len(tasks))
	for i := range tasks {
		item := &tasks[i]
		events = append(events, taskEvent{line: item.LineNumber, task: item, depth: item.Depth})
		for j := range item.Comments {
			comment := &item.Comments[j]
			events = append(events, taskEvent{line: comment.LineNumber, task: item, comment: comment, depth: item.Depth + 1})
		}
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].line < events[j].line })

	lastSection := ""
	for _, event := range events {
		if event.comment == nil {
			sectionPath := strings.Join(event.task.Sections, " / ")
			if sectionPath != lastSection {
				if lastSection != "" {
					fmt.Fprintln(output)
				}
				if sectionPath != "" {
					fmt.Fprintf(output, "%s:\n", sectionPath)
				}
				lastSection = sectionPath
			}
			checkbox := "[ ]"
			if event.task.Done {
				checkbox = "[x]"
			}
			fmt.Fprintf(output, "%s- %s %s", strings.Repeat("  ", event.depth), checkbox, event.task.Description)
			for _, tag := range event.task.Tags {
				fmt.Fprintf(output, " #%s", tag)
			}
			if event.task.Status != "" {
				fmt.Fprintf(output, " @%s", event.task.Status)
			}
			for _, project := range event.task.Projects {
				fmt.Fprintf(output, " +%s", project)
			}
			if number, ok := localTaskNumber(event.task.ID, project); ok {
				fmt.Fprintf(output, " $%s", number)
			} else if event.task.ID != "" {
				fmt.Fprintf(output, " id:%s", event.task.ID)
			}
			if event.task.UpdatedAt != "" {
				timestamp, err := time.Parse(time.RFC3339Nano, event.task.UpdatedAt)
				if err == nil {
					fmt.Fprintf(output, " %s", task.FormatTimestamp(timestamp))
				}
			}
			fmt.Fprintln(output)
			continue
		}
		fmt.Fprintf(output, "%s- %s", strings.Repeat("  ", event.depth), event.comment.Text)
		for _, tag := range event.comment.Tags {
			fmt.Fprintf(output, " #%s", tag)
		}
		for _, project := range event.comment.Projects {
			fmt.Fprintf(output, " +%s", project)
		}
		fmt.Fprintln(output)
	}
}
