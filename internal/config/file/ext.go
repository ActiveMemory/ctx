//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

// File extension constants.
const (
	// ExtMarkdown is the Markdown file extension.
	ExtMarkdown = ".md"
	// ExtTxt is the plain text file extension.
	ExtTxt = ".txt"
	// ExtGo is the Go source file extension.
	ExtGo = ".go"
	// ExtJSONL is the JSON Lines file extension.
	ExtJSONL = ".jsonl"
	// ExtYAML is the YAML file extension.
	ExtYAML = ".yaml"
	// ExtJSON is the JSON file extension.
	ExtJSON = ".json"
	// ExtEnc is the encrypted file extension.
	ExtEnc = ".enc"
	// ExtSh is the shell script file extension.
	ExtSh = ".sh"
	// ExtPs1 is the PowerShell script file extension.
	ExtPs1 = ".ps1"
	// ExtTmp is the temporary file suffix for atomic writes.
	ExtTmp = ".tmp"
	// ExtExample is the suffix for example/template files that are safe
	// to have in the working directory (e.g., .env.example).
	ExtExample = ".example"
	// ExtSample is the suffix for sample files that are safe to have
	// in the working directory (e.g., config.sample).
	ExtSample = ".sample"
)

// BackupFormat is the format string for timestamped backup file names.
// Args: original filename, Unix timestamp.
const BackupFormat = "%s.%d.bak"
