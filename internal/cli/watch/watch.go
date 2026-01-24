//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/cli/init"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	watchLog      string
	watchDryRun   bool
	watchAutoSave bool
)

// WatchCmd returns the watch command.
func WatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for context-update commands in AI output",
		Long: `Watch stdin or a log file for <context-update> commands and apply them.

This command parses AI output looking for structured update commands:

  <context-update type="task">Implement user auth</context-update>
  <context-update type="decision">Use PostgreSQL</context-update>
  <context-update type="learning">Mock functions must be hoisted</context-update>
  <context-update type="complete">user auth</context-update>

Use --log to watch a specific file instead of stdin.
Use --dry-run to see what would be updated without making changes.
Use --auto-save to periodically save session snapshots (every 5 updates).

Press Ctrl+C to stop watching.`,
		RunE: runWatch,
	}

	cmd.Flags().StringVar(&watchLog, "log", "", "Log file to watch (default: stdin)")
	cmd.Flags().BoolVar(&watchDryRun, "dry-run", false, "Show updates without applying")
	cmd.Flags().BoolVar(&watchAutoSave, "auto-save", false, "Save session snapshots periodically")

	return cmd
}

// ContextUpdate represents a parsed context update command.
type ContextUpdate struct {
	Type    string
	Content string
}

func runWatch(cmd *cobra.Command, _ []string) error {
	// Check if context exists
	if !context.Exists("") {
		return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	cmd.Println(cyan("Watching for context updates..."))
	if watchDryRun {
		yellow := color.New(color.FgYellow).SprintFunc()
		cmd.Println(yellow("DRY RUN — No changes will be made"))
	}
	cmd.Println("Press Ctrl+C to stop")
	cmd.Println()

	var reader io.Reader
	if watchLog != "" {
		file, err := os.Open(watchLog)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("failed to close log file: %v\n", err)
			}
		}(file)
		reader = file
	} else {
		reader = os.Stdin
	}

	return processStream(cmd, reader)
}

// autoSaveInterval is the number of updates between auto-saves.
const autoSaveInterval = 5

func processStream(cmd *cobra.Command, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	// Use a larger buffer for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Pattern to match context-update tags
	updatePattern := regexp.MustCompile(`<context-update\s+type="([^"]+)"[^>]*>([^<]+)</context-update>`)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Track applied updates for auto-save
	updateCount := 0
	var appliedUpdates []ContextUpdate

	for scanner.Scan() {
		line := scanner.Text()

		// Check for context-update commands
		matches := updatePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				update := ContextUpdate{
					Type:    strings.ToLower(match[1]),
					Content: strings.TrimSpace(match[2]),
				}

				if watchDryRun {
					cmd.Printf("%s Would apply: [%s] %s\n", yellow("○"), update.Type, update.Content)
				} else {
					err := applyUpdate(update)
					if err != nil {
						cmd.Printf("%s Failed to apply [%s]: %v\n", color.RedString("✗"), update.Type, err)
					} else {
						cmd.Printf("%s Applied: [%s] %s\n", green("✓"), update.Type, update.Content)
						updateCount++
						appliedUpdates = append(appliedUpdates, update)

						// Auto-save every N updates
						if watchAutoSave && updateCount%autoSaveInterval == 0 {
							if err := watchAutoSaveSession(appliedUpdates); err != nil {
								cmd.Printf("%s Auto-save failed: %v\n", yellow("⚠"), err)
							} else {
								cmd.Printf("%s Auto-saved session after %d updates\n", cyan("📸"), updateCount)
							}
						}
					}
				}
			}
		}
	}

	// Final auto-save if there are remaining updates
	if watchAutoSave && len(appliedUpdates) > 0 && updateCount%autoSaveInterval != 0 {
		if err := watchAutoSaveSession(appliedUpdates); err != nil {
			cmd.Printf("%s Final auto-save failed: %v\n", yellow("⚠"), err)
		} else {
			cmd.Printf("%s Final auto-save completed (%d total updates)\n", cyan("📸"), updateCount)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func applyUpdate(update ContextUpdate) error {
	switch update.Type {
	case "task":
		return applyTaskUpdate(update.Content)
	case "decision":
		return applyDecisionUpdate(update.Content)
	case "learning":
		return applyLearningUpdate(update.Content)
	case "convention":
		return applyConventionUpdate(update.Content)
	case "complete":
		return applyCompleteUpdate(update.Content)
	default:
		return fmt.Errorf("unknown update type: %s", update.Type)
	}
}

func applyTaskUpdate(content string) error {
	// Reuse the add command logic
	args := []string{"task", content}
	return add.runAdd(nil, args)
}

func applyDecisionUpdate(content string) error {
	args := []string{"decision", content}
	// Suppress output from add command during watch
	return runAddSilent(args)
}

func applyLearningUpdate(content string) error {
	args := []string{"learning", content}
	return runAddSilent(args)
}

func applyConventionUpdate(content string) error {
	args := []string{"convention", content}
	return runAddSilent(args)
}

func applyCompleteUpdate(content string) error {
	args := []string{content}
	return runCompleteSilent(args)
}

// runAddSilent runs the add command without output
func runAddSilent(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}

	fileType := strings.ToLower(args[0])
	content := strings.Join(args[1:], " ")

	fileName, ok := add.fileTypeMap[fileType]
	if !ok {
		return fmt.Errorf("unknown type %q", fileType)
	}

	filePath := init.contextDirName + "/" + fileName

	existing, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var entry string
	switch fileType {
	case "decision", "decisions":
		entry = add.formatDecision(content)
	case "task", "tasks":
		entry = add.formatTask(content, "")
	case "learning", "learnings":
		entry = add.formatLearning(content)
	case "convention", "conventions":
		entry = add.formatConvention(content)
	}

	newContent := add.appendEntry(existing, entry, fileType, "")
	return os.WriteFile(filePath, newContent, 0644)
}

// runCompleteSilent runs the complete command without output
func runCompleteSilent(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no task specified")
	}

	query := args[0]
	filePath := init.contextDirName + "/TASKS.md"

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	taskPattern := regexp.MustCompile(`^(\s*)-\s*\[\s*]\s*(.+)$`)

	matchedLine := -1
	for i, line := range lines {
		matches := taskPattern.FindStringSubmatch(line)
		if matches != nil {
			taskText := matches[2]
			if strings.Contains(strings.ToLower(taskText), strings.ToLower(query)) {
				matchedLine = i
				break
			}
		}
	}

	if matchedLine == -1 {
		return fmt.Errorf("no task matching %q found", query)
	}

	lines[matchedLine] = taskPattern.ReplaceAllString(lines[matchedLine], "$1- [x] $2")
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// watchAutoSaveSession saves a session snapshot during watch mode.
func watchAutoSaveSession(updates []ContextUpdate) error {
	// Ensure sessions directory exists
	sessionsDir := filepath.Join(init.contextDirName, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-watch.md", now.Format("2006-01-02-150405"))
	filePath := filepath.Join(sessionsDir, filename)

	// Build session content
	content := buildWatchSession(now, updates)

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// buildWatchSession creates a session snapshot from watch mode updates.
func buildWatchSession(timestamp time.Time, updates []ContextUpdate) string {
	var sb strings.Builder

	sb.WriteString("# Watch Mode Session\n\n")
	sb.WriteString(fmt.Sprintf("**Date**: %s\n", timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s\n", timestamp.Format("15:04:05")))
	sb.WriteString("**Type**: watch-auto-save\n\n")
	sb.WriteString("---\n\n")

	sb.WriteString("## Applied Updates\n\n")

	// Group updates by type
	updatesByType := make(map[string][]string)
	for _, u := range updates {
		updatesByType[u.Type] = append(updatesByType[u.Type], u.Content)
	}

	// Write updates by type
	typeOrder := []string{"task", "decision", "learning", "convention", "complete"}
	for _, t := range typeOrder {
		contents, ok := updatesByType[t]
		if !ok || len(contents) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("### %ss\n\n", strings.ToUpper(t[:1])+t[1:]))
		for _, c := range contents {
			sb.WriteString(fmt.Sprintf("- %s\n", c))
		}
		sb.WriteString("\n")
	}

	// Add the current context snapshot
	sb.WriteString("---\n\n")
	sb.WriteString("## Context Snapshot\n\n")

	// Read TASKS.md
	tasksPath := filepath.Join(init.contextDirName, "TASKS.md")
	if tasksContent, err := os.ReadFile(tasksPath); err == nil {
		sb.WriteString("### Current Tasks\n\n")
		sb.WriteString("```markdown\n")
		sb.WriteString(string(tasksContent))
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}
