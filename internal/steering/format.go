//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"bytes"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// syncableTool returns true if the tool supports native-format sync.
func syncableTool(tool string) bool {
	for _, t := range syncableTools {
		if t == tool {
			return true
		}
	}
	return false
}

// nativePath returns the output file path for a steering file in the
// given tool's native format.
func nativePath(
	projectRoot, tool, name string,
) string {
	switch tool {
	case cfgHook.ToolCursor:
		return filepath.Join(
			projectRoot, cfgSteering.DirCursorDot,
			cfgSteering.DirRules, name+cfgSteering.ExtMDC,
		)
	case cfgHook.ToolCline:
		return filepath.Join(
			projectRoot, cfgSteering.DirClinerules,
			name+file.ExtMarkdown,
		)
	case cfgHook.ToolKiro:
		return filepath.Join(
			projectRoot, cfgSteering.DirKiroDot,
			cfgSteering.DirSteering, name+file.ExtMarkdown,
		)
	default:
		return ""
	}
}

// validateOutputPath checks that the output path resolves within the
// project root boundary. This prevents path traversal via crafted
// steering file names.
func validateOutputPath(outPath, projectRoot string) error {
	absOut, absOutErr := filepath.Abs(outPath)
	if absOutErr != nil {
		return errSteering.ResolveOutput(absOutErr)
	}
	absRoot, absRootErr := filepath.Abs(projectRoot)
	if absRootErr != nil {
		return errSteering.ResolveRoot(absRootErr)
	}

	rel, relErr := filepath.Rel(absRoot, absOut)
	if relErr != nil {
		return errSteering.ComputeRelPath(relErr)
	}

	// Reject paths that escape the project root.
	escape := token.ParentDir + string(filepath.Separator)
	if strings.HasPrefix(rel, escape) || rel == token.ParentDir {
		return errSteering.OutputEscapesRoot(outPath, projectRoot)
	}

	return nil
}

// formatNative converts a steering file to the tool's native format.
func formatNative(tool string, sf *SteeringFile) []byte {
	switch tool {
	case cfgHook.ToolCursor:
		return formatCursor(sf)
	case cfgHook.ToolCline:
		return formatCline(sf)
	case cfgHook.ToolKiro:
		return formatKiro(sf)
	default:
		return nil
	}
}

// formatCursor produces Cursor-compatible .mdc content with frontmatter.
func formatCursor(sf *SteeringFile) []byte {
	fm := cursorFrontmatter{
		Description: sf.Description,
		Globs:       []any{},
		AlwaysApply: sf.Inclusion == cfgSteering.InclusionAlways,
	}

	raw, _ := yaml.Marshal(fm)

	var buf bytes.Buffer
	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])
	buf.Write(raw)
	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])
	if sf.Body != "" {
		buf.WriteString(sf.Body)
	}
	return buf.Bytes()
}

// formatCline produces Cline-compatible plain markdown (no frontmatter).
func formatCline(sf *SteeringFile) []byte {
	var buf bytes.Buffer
	buf.WriteString(token.HeadingLevelOneStart)
	buf.WriteString(sf.Name)
	buf.WriteString(token.DoubleNewline)
	if sf.Body != "" {
		buf.WriteString(sf.Body)
	}
	return buf.Bytes()
}

// formatKiro produces Kiro-compatible steering file with frontmatter.
func formatKiro(sf *SteeringFile) []byte {
	fm := kiroFrontmatter{
		Name:        sf.Name,
		Description: sf.Description,
		Mode:        mapKiroMode(sf.Inclusion),
	}

	raw, _ := yaml.Marshal(fm)

	var buf bytes.Buffer
	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])
	buf.Write(raw)
	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])
	if sf.Body != "" {
		buf.WriteString(sf.Body)
	}
	return buf.Bytes()
}

// mapKiroMode maps ctx inclusion modes to Kiro equivalents.
func mapKiroMode(
	inc cfgSteering.InclusionMode,
) string {
	switch inc {
	case cfgSteering.InclusionAlways:
		return string(cfgSteering.InclusionAlways)
	case cfgSteering.InclusionAuto:
		return string(cfgSteering.InclusionAuto)
	case cfgSteering.InclusionManual:
		return string(cfgSteering.InclusionManual)
	default:
		return string(cfgSteering.InclusionManual)
	}
}

// unchanged returns true if the file at path already exists and has
// the same content as data.
func unchanged(path string, data []byte) bool {
	existing, err := ctxIo.SafeReadUserFile(path)
	if err != nil {
		return false
	}
	return bytes.Equal(existing, data)
}

// writeFile creates parent directories as needed and writes data to path.
func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if mkdirErr := ctxIo.SafeMkdirAll(dir, fs.PermExec); mkdirErr != nil {
		return mkdirErr
	}
	return ctxIo.SafeWriteFile(path, data, fs.PermFile)
}
