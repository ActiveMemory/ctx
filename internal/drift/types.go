package drift

// IssueType categorizes a drift issue for grouping and filtering.
type IssueType string

const (
	// IssueDeadPath indicates a file path reference that no longer exists.
	IssueDeadPath IssueType = "dead_path"
	// IssueStaleness indicates accumulated completed tasks needing archival.
	IssueStaleness IssueType = "staleness"
	// IssueSecret indicates a file that may contain secrets or credentials.
	IssueSecret IssueType = "potential_secret"
	// IssueMissing indicates a required context file that does not exist.
	IssueMissing IssueType = "missing_file"
	// IssueStaleAge indicates a context file that hasn't been modified recently.
	IssueStaleAge IssueType = "stale_age"
	// IssueEntryCount indicates a knowledge file has too many entries.
	IssueEntryCount IssueType = "entry_count"
	// IssueMissingPackage indicates an internal package not documented in ARCHITECTURE.md.
	IssueMissingPackage IssueType = "missing_package"
)

// StatusType represents the overall status of a drift report.
type StatusType string

const (
	// StatusOk means no drift was detected.
	StatusOk StatusType = "ok"
	// StatusWarning means non-critical issues were found.
	StatusWarning StatusType = "warning"
	// StatusViolation means constitution violations were found.
	StatusViolation StatusType = "violation"
)

// CheckName identifies a drift detection check.
type CheckName string

const (
	// CheckPathReferences validates that file paths in context files exist.
	CheckPathReferences CheckName = "path_references"
	// CheckStaleness detects accumulated completed tasks.
	CheckStaleness CheckName = "staleness_check"
	// CheckConstitution verifies constitution rules are respected.
	CheckConstitution CheckName = "constitution_check"
	// CheckRequiredFiles ensures all required context files are present.
	CheckRequiredFiles CheckName = "required_files"
	// CheckFileAge checks whether context files have been modified recently.
	CheckFileAge CheckName = "file_age_check"
	// CheckEntryCount checks whether knowledge files have excessive entries.
	CheckEntryCount CheckName = "entry_count_check"
	// CheckMissingPackages checks for undocumented internal packages.
	CheckMissingPackages CheckName = "missing_packages"
)

// Issue represents a detected drift issue.
//
// Issues are categorized by type and may reference specific files, lines,
// or paths in the codebase.
//
// Fields:
//   - File: Context file where the issue was detected (e.g., "ARCHITECTURE.md")
//   - Line: Line number in the file, if applicable
//   - Type: Issue category (e.g., "dead_path", "staleness", "missing_file")
//   - Message: Human-readable description of the issue
//   - Path: Referenced path that caused the issue, if applicable
//   - Rule: Constitution rule that was violated, if applicable
type Issue struct {
	File    string    `json:"file"`
	Line    int       `json:"line,omitempty"`
	Type    IssueType `json:"type"`
	Message string    `json:"message"`
	Path    string    `json:"path,omitempty"`
	Rule    string    `json:"rule,omitempty"`
}

// Report represents the complete drift detection report.
//
// Contains categorized issues and a list of checks that passed.
//
// Fields:
//   - Warnings: Non-critical issues that should be addressed
//   - Violations: Critical issues that indicate constitution violations
//   - Passed: Names of checks that are completed without issues
type Report struct {
	Warnings   []Issue     `json:"warnings"`
	Violations []Issue     `json:"violations"`
	Passed     []CheckName `json:"passed"`
}
