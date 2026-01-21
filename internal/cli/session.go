package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	sessionsDirName = ".context/sessions"
)

// SessionCmd returns the session command with subcommands.
func SessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage session snapshots",
		Long: `Manage session snapshots in .context/sessions/.

Sessions capture the state of your context at a point in time,
including current tasks, recent decisions, and learnings.

Subcommands:
  save    Save current context state to a session file
  list    List saved sessions with summaries
  load    Load and display a previous session
  parse   Convert .jsonl transcript to readable markdown`,
	}

	cmd.AddCommand(sessionSaveCmd())
	cmd.AddCommand(sessionListCmd())

	return cmd
}

var (
	sessionTopic string
	sessionType  string
)

// sessionSaveCmd returns the session save subcommand.
func sessionSaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save [topic]",
		Short: "Save current context state to a session file",
		Long: `Save a snapshot of the current context state to .context/sessions/.

The session file includes:
  - Summary of what was done
  - Current tasks from TASKS.md
  - Recent decisions from DECISIONS.md
  - Recent learnings from LEARNINGS.md

Examples:
  ctx session save "implemented auth"
  ctx session save "refactored API" --type feature
  ctx session save  # prompts for topic interactively`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSessionSave,
	}

	cmd.Flags().StringVarP(&sessionType, "type", "t", "session", "Session type (feature, bugfix, refactor, session)")

	return cmd
}

func runSessionSave(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Get topic from args or use default
	topic := "manual-save"
	if len(args) > 0 {
		topic = args[0]
	}

	// Sanitize topic for filename
	topic = sanitizeFilename(topic)

	// Ensure sessions directory exists
	if err := os.MkdirAll(sessionsDirName, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate filename
	now := time.Now()
	filename := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02-150405"), topic)
	filePath := filepath.Join(sessionsDirName, filename)

	// Build session content
	content, err := buildSessionContent(topic, sessionType, now)
	if err != nil {
		return fmt.Errorf("failed to build session content: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	fmt.Printf("%s Session saved to %s\n", green("✓"), filePath)
	return nil
}

var listLimit int

// sessionListCmd returns the session list subcommand.
func sessionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved sessions with summaries",
		Long: `List all saved sessions in .context/sessions/.

Shows session date, topic, type, and a brief summary for each session.
Sessions are sorted by date (newest first).

Examples:
  ctx session list
  ctx session list --limit 5`,
		RunE: runSessionList,
	}

	cmd.Flags().IntVarP(&listLimit, "limit", "n", 10, "Maximum number of sessions to display")

	return cmd
}

func runSessionList(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDirName); os.IsNotExist(err) {
		fmt.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(sessionsDirName)
	if err != nil {
		return fmt.Errorf("failed to read sessions directory: %w", err)
	}

	// Filter and collect session files
	var sessions []sessionInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only show .md files (not .jsonl transcripts)
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		// Skip summary files that accompany jsonl files
		if strings.HasSuffix(name, "-summary.md") {
			continue
		}

		info, err := parseSessionFile(filepath.Join(sessionsDirName, name))
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		info.Filename = name
		sessions = append(sessions, info)
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found. Use 'ctx session save' to create one.")
		return nil
	}

	// Sort by date (newest first) - filenames are date-prefixed so reverse sort works
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	// Limit output
	if listLimit > 0 && len(sessions) > listLimit {
		sessions = sessions[:listLimit]
	}

	// Display
	fmt.Printf("Sessions in %s:\n\n", sessionsDirName)
	for _, s := range sessions {
		fmt.Printf("%s %s\n", cyan("●"), s.Topic)
		fmt.Printf("  %s %s | %s %s\n",
			gray("Date:"), s.Date,
			gray("Type:"), s.Type)
		if s.Summary != "" {
			fmt.Printf("  %s %s\n", gray("Summary:"), truncate(s.Summary, 60))
		}
		fmt.Printf("  %s %s\n", yellow("File:"), s.Filename)
		fmt.Println()
	}

	fmt.Printf("Total: %d session(s)\n", len(sessions))
	return nil
}

// sessionInfo holds parsed information about a session file.
type sessionInfo struct {
	Filename string
	Topic    string
	Date     string
	Type     string
	Summary  string
}

// parseSessionFile extracts metadata from a session file.
func parseSessionFile(path string) (sessionInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return sessionInfo{}, err
	}

	contentStr := string(content)
	info := sessionInfo{}

	// Extract topic from first line (# Session: topic)
	if strings.HasPrefix(contentStr, "# Session:") {
		lineEnd := strings.Index(contentStr, "\n")
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(contentStr[11:lineEnd])
		}
	} else if strings.HasPrefix(contentStr, "# ") {
		// Alternative format: # Topic
		lineEnd := strings.Index(contentStr, "\n")
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(contentStr[2:lineEnd])
		}
	}

	// Extract date
	if idx := strings.Index(contentStr, "**Date**:"); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], "\n")
		if lineEnd != -1 {
			info.Date = strings.TrimSpace(contentStr[idx+9 : idx+lineEnd])
		}
	}

	// Extract type
	if idx := strings.Index(contentStr, "**Type**:"); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], "\n")
		if lineEnd != -1 {
			info.Type = strings.TrimSpace(contentStr[idx+9 : idx+lineEnd])
		}
	}

	// Extract summary (first non-empty line after ## Summary)
	if idx := strings.Index(contentStr, "## Summary"); idx != -1 {
		afterSummary := contentStr[idx+10:]
		lines := strings.Split(afterSummary, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "[") {
				info.Summary = line
				break
			}
		}
	}

	return info, nil
}

// truncate shortens a string to maxLen characters, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// sanitizeFilename converts a topic string to a safe filename component.
func sanitizeFilename(s string) string {
	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	s = re.ReplaceAllString(s, "-")
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	// Convert to lowercase
	s = strings.ToLower(s)
	// Limit length
	if len(s) > 50 {
		s = s[:50]
	}
	if s == "" {
		s = "session"
	}
	return s
}

// buildSessionContent creates the markdown content for a session file.
func buildSessionContent(topic, sessionType string, timestamp time.Time) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Session: %s\n\n", topic))
	sb.WriteString(fmt.Sprintf("**Date**: %s\n", timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s\n", timestamp.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("**Type**: %s\n", sessionType))
	sb.WriteString("\n---\n\n")

	// Summary section (placeholder for user to fill in)
	sb.WriteString("## Summary\n\n")
	sb.WriteString("[Describe what was accomplished in this session]\n\n")
	sb.WriteString("---\n\n")

	// Current Tasks
	sb.WriteString("## Current Tasks\n\n")
	tasks, err := readContextSection("TASKS.md", "## In Progress", "## Next Up")
	if err == nil && tasks != "" {
		sb.WriteString("### In Progress\n\n")
		sb.WriteString(tasks)
		sb.WriteString("\n")
	}
	nextTasks, err := readContextSection("TASKS.md", "## Next Up", "## Completed")
	if err == nil && nextTasks != "" {
		sb.WriteString("### Next Up\n\n")
		sb.WriteString(nextTasks)
		sb.WriteString("\n")
	}
	sb.WriteString("---\n\n")

	// Recent Decisions
	sb.WriteString("## Recent Decisions\n\n")
	decisions, err := readRecentDecisions()
	if err == nil && decisions != "" {
		sb.WriteString(decisions)
	} else {
		sb.WriteString("[No recent decisions found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Recent Learnings
	sb.WriteString("## Recent Learnings\n\n")
	learnings, err := readRecentLearnings()
	if err == nil && learnings != "" {
		sb.WriteString(learnings)
	} else {
		sb.WriteString("[No recent learnings found]\n")
	}
	sb.WriteString("\n---\n\n")

	// Tasks for Next Session
	sb.WriteString("## Tasks for Next Session\n\n")
	sb.WriteString("[List tasks to continue in the next session]\n\n")

	return sb.String(), nil
}

// readContextSection reads a section from a context file between two headers.
func readContextSection(filename, startHeader, endHeader string) (string, error) {
	filePath := filepath.Join(contextDirName, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find start
	startIdx := strings.Index(contentStr, startHeader)
	if startIdx == -1 {
		return "", fmt.Errorf("section not found: %s", startHeader)
	}
	startIdx += len(startHeader)

	// Find end
	endIdx := len(contentStr)
	if endHeader != "" {
		idx := strings.Index(contentStr[startIdx:], endHeader)
		if idx != -1 {
			endIdx = startIdx + idx
		}
	}

	section := strings.TrimSpace(contentStr[startIdx:endIdx])
	return section, nil
}

// readRecentDecisions extracts the most recent decisions from DECISIONS.md.
func readRecentDecisions() (string, error) {
	filePath := filepath.Join(contextDirName, "DECISIONS.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find decision headers (## [YYYY-MM-DD] Title)
	re := regexp.MustCompile(`(?m)^## \[\d{4}-\d{2}-\d{2}\].*$`)
	matches := re.FindAllStringIndex(contentStr, -1)

	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 3 decisions (most recent)
	limit := 3
	if len(matches) < limit {
		limit = len(matches)
	}

	var decisions []string
	for i := len(matches) - limit; i < len(matches); i++ {
		start := matches[i][0]
		end := len(contentStr)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		decision := strings.TrimSpace(contentStr[start:end])
		// Only include the header for brevity
		headerEnd := strings.Index(decision, "\n")
		if headerEnd != -1 {
			decisions = append(decisions, "- "+decision[:headerEnd])
		}
	}

	return strings.Join(decisions, "\n"), nil
}

// readRecentLearnings extracts the most recent learnings from LEARNINGS.md.
func readRecentLearnings() (string, error) {
	filePath := filepath.Join(contextDirName, "LEARNINGS.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Find learning entries (- **[YYYY-MM-DD]** text)
	re := regexp.MustCompile(`(?m)^- \*\*\[\d{4}-\d{2}-\d{2}\]\*\*.*$`)
	matches := re.FindAllString(contentStr, -1)

	if len(matches) == 0 {
		return "", nil
	}

	// Get the last 5 learnings (most recent)
	limit := 5
	if len(matches) < limit {
		limit = len(matches)
	}

	return strings.Join(matches[len(matches)-limit:], "\n"), nil
}
