//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

// Packages maps manifest filenames to their ecosystem descriptions.
// Used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}

// Go toolchain command constants.
const (
	// GoBinary is the Go executable name.
	GoBinary = "go"
	// GoList is the go list subcommand.
	GoList = "list"
	// GoFlagJSON is the JSON output flag for go list.
	GoFlagJSON = "-json"
	// GoAllPackages is the recursive package pattern.
	GoAllPackages = "./..."
)

// Cargo toolchain command constants.
const (
	// CargoBinary is the Cargo executable name.
	CargoBinary = "cargo"
	// CargoMetadataCmd is the cargo metadata subcommand.
	CargoMetadataCmd = "metadata"
	// CargoFlagFormatVersion is the format version flag.
	CargoFlagFormatVersion = "--format-version"
	// CargoFormatVersion1 is the format version value.
	CargoFormatVersion1 = "1"
	// CargoFlagNoDeps skips dependency resolution.
	CargoFlagNoDeps = "--no-deps"
)

// Python manifest files.
const (
	// FileRequirements is the requirements.txt manifest.
	FileRequirements = "requirements.txt"
	// FilePyproject is the pyproject.toml manifest.
	FilePyproject = "pyproject.toml"
)

// Python dependency section keys for pyproject.toml parsing.
const (
	// PyDeps is the main dependencies section suffix.
	PyDeps = "dependencies"
	// PyDevDeps is the dev-dependencies section suffix.
	PyDevDeps = "dev-dependencies"
	// PyDev is the short dev section suffix.
	PyDev = "dev"
	// PyGraphRoot is the root node label in Python dependency graphs.
	PyGraphRoot = "project"
)

// PyMetaKeys lists pyproject.toml keys that are metadata, not dependencies.
// Entries matching these keys are skipped during Poetry-style parsing.
var PyMetaKeys = map[string]bool{
	"name":        true,
	"version":     true,
	"description": true,
}

// PyVersionSpecChars contains the characters that begin a version
// specifier or extras bracket in Python dependency lines.
const PyVersionSpecChars = "><=!~["

// TOML parsing tokens.
const (
	// TomlComment is the inline comment prefix.
	TomlComment = " #"
	// TomlOptionPrefix is the prefix for pip options in requirements.txt.
	TomlOptionPrefix = "-"
	// TomlSectionOpen is the opening bracket for TOML sections.
	TomlSectionOpen = "["
	// TomlArrayAssign variants for matching inline arrays.
	TomlArrayAssign1 = " = ["
	TomlArrayAssign2 = "= ["
	TomlArrayAssign3 = " =["
	TomlArrayAssign4 = "=["
	// TomlSemicolon is the environment marker separator.
	TomlSemicolon = ";"
	// TomlProjectPrefix is the TOML section path for project tables.
	TomlProjectPrefix = "project."
	// TomlPoetryPrefix is the TOML section path for Poetry tables.
	TomlPoetryPrefix = "tool.poetry."
	// TomlSectionFmt is the Printf format for a TOML section header.
	// Args: prefix (e.g. "project."), suffix (e.g. "dependencies").
	TomlSectionFmt = "[%s%s]"
)

// Table formatting constants for dependency output.
const (
	// TableColPackage is the column width for package names in table output.
	TableColPackage = 50
	// TableColImports is the column width for import lists in table output.
	TableColImports = 30
	// TableHeaderPackage is the column header for package names.
	TableHeaderPackage = "Package"
	// TableHeaderImports is the column header for import lists.
	TableHeaderImports = "Imports"
	// TableRowFormat is the dynamic-width row format template for table output.
	TableRowFormat = "%%-%ds %%s\n"
	// MermaidEdgeFormat is the Mermaid graph edge format string.
	MermaidEdgeFormat = "    %s[\"%s\"] --> %s[\"%s\"]\n"
)
