//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// agyPluginManifest is the subset of internal/assets/agy/plugin.json
// inspected by the guard.
type agyPluginManifest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

// TestAgyPluginManifest asserts that the Antigravity plugin manifest
// has name == "ctx", carries the repo version, points to a valid logo
// that resolves in the embedded FS, and that its rules/AGENTS.md equals
// integrations/agents.md.
func TestAgyPluginManifest(t *testing.T) {
	data, err := FS.ReadFile(asset.PathAgyPluginJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathAgyPluginJSON, err)
	}
	var m agyPluginManifest
	if unmarshalErr := json.Unmarshal(data, &m); unmarshalErr != nil {
		t.Fatalf("parse %s: %v", asset.PathAgyPluginJSON, unmarshalErr)
	}
	if m.Name != "ctx" {
		t.Errorf("name = %q, want %q", m.Name, "ctx")
	}
	if want := repoVersion(t); m.Version != want {
		t.Errorf("version = %q, want %q (VERSION file)", m.Version, want)
	}
	if m.Logo == "" {
		t.Error("logo is empty")
	} else {
		logoPath := path.Join(asset.DirAgy, m.Logo)
		if _, statErr := fs.Stat(FS, logoPath); statErr != nil {
			t.Errorf(
				"logo %q not found in embedded FS (%s): %v",
				m.Logo, logoPath, statErr,
			)
		}
	}

	rulesData, err := FS.ReadFile(asset.PathAgyRulesAgentsMd)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathAgyRulesAgentsMd, err)
	}
	integrationsData, err := FS.ReadFile(asset.PathAgentsMd)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathAgentsMd, err)
	}
	if !bytes.Equal(rulesData, integrationsData) {
		t.Errorf(
			"%s does not match %s",
			asset.PathAgyRulesAgentsMd, asset.PathAgentsMd,
		)
	}
}

// agyMCPConfig represents the structure of agy/mcp_config.json.
type agyMCPConfig struct {
	MCPServers map[string]struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"mcpServers"`
}

// TestAgyMCPConfig asserts that mcp_config.json registers the ctx MCP
// server with command "ctx" and args ["mcp", "serve"].
func TestAgyMCPConfig(t *testing.T) {
	data, err := FS.ReadFile(asset.PathAgyMCPJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathAgyMCPJSON, err)
	}
	var cfg agyMCPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse %s: %v", asset.PathAgyMCPJSON, err)
	}
	server, ok := cfg.MCPServers["ctx"]
	if !ok {
		t.Fatalf("%s missing mcpServers.ctx", asset.PathAgyMCPJSON)
	}
	if server.Command != "ctx" {
		t.Errorf("command = %q, want %q", server.Command, "ctx")
	}
	wantArgs := []string{"mcp", "serve"}
	if !slices.Equal(server.Args, wantArgs) {
		t.Errorf("args = %v, want %v", server.Args, wantArgs)
	}
}

// TestAgyHooksManifestShape asserts that hooks.json has the top-level
// "ctx" group, valid lifecycle events, and non-empty commands.
func TestAgyHooksManifestShape(t *testing.T) {
	data, err := FS.ReadFile(asset.PathAgyHooksJSON)
	if err != nil {
		t.Fatalf("read %s: %v", asset.PathAgyHooksJSON, err)
	}
	var manifest map[string]map[string]json.RawMessage
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse %s: %v", asset.PathAgyHooksJSON, err)
	}
	ctxHooks, ok := manifest["ctx"]
	if !ok {
		t.Fatalf(
			"%s missing top-level %q hook group",
			asset.PathAgyHooksJSON, "ctx",
		)
	}

	validEvents := []string{
		"SessionStart",
		"PreInvocation",
		"PreToolUse",
		"PostToolUse",
		"PostInvocation",
		"Stop",
	}

	for event := range ctxHooks {
		if !slices.Contains(validEvents, event) {
			t.Errorf("unexpected event %q in hooks.json", event)
		}
	}
	for _, event := range validEvents {
		raw, ok := ctxHooks[event]
		if !ok {
			t.Errorf("missing expected event %q in hooks.json", event)
			continue
		}
		var flatHandlers []struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		}
		var matcherGroups []struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Type    string `json:"type"`
				Command string `json:"command"`
			} `json:"hooks"`
		}

		if event == "PreToolUse" || event == "PostToolUse" {
			if err := json.Unmarshal(raw, &matcherGroups); err != nil {
				t.Errorf("%s: unmarshal matcher groups: %v", event, err)
			}
			if len(matcherGroups) == 0 {
				t.Errorf("%s: no matcher groups configured", event)
			}
			for i, g := range matcherGroups {
				if g.Matcher == "" {
					t.Errorf("%s group %d: empty matcher", event, i)
				}
				if len(g.Hooks) == 0 {
					t.Errorf("%s group %d: no hooks configured", event, i)
				}
				for j, h := range g.Hooks {
					if h.Type != "command" {
						t.Errorf(
							"%s group %d hook %d: type = %q, want %q",
							event, i, j, h.Type, "command",
						)
					}
					if strings.TrimSpace(h.Command) == "" {
						t.Errorf(
							"%s group %d hook %d: command is empty",
							event, i, j,
						)
					}
				}
			}
		} else {
			if err := json.Unmarshal(raw, &flatHandlers); err != nil {
				t.Errorf("%s: unmarshal flat handlers: %v", event, err)
			}
			if len(flatHandlers) == 0 {
				t.Errorf("%s: no handlers configured", event)
			}
			for i, h := range flatHandlers {
				if h.Type != "command" {
					t.Errorf(
						"%s handler %d: type = %q, want %q",
						event, i, h.Type, "command",
					)
				}
				if strings.TrimSpace(h.Command) == "" {
					t.Errorf("%s handler %d: command is empty", event, i)
				}
			}
		}
	}
}

// TestAgySkills asserts that Antigravity skills match the Codex skills list,
// contain no `allowed-tools:` frontmatter, stay under the 24,000-byte body
// cap, and all cited reference files exist in the embedded FS.
func TestAgySkills(t *testing.T) {
	codexSkills := skillDirs(t, asset.DirCodexSkills)
	agySkills := skillDirs(t, asset.DirAgySkills)

	if !slices.Equal(agySkills, codexSkills) {
		t.Errorf(
			"Antigravity skills do not match Codex skills\n agy: %v\n codex: %v",
			agySkills, codexSkills,
		)
	}

	reRef := regexp.MustCompile(`references/[A-Za-z0-9._-]+`)
	fence := []byte("---")

	for _, name := range agySkills {
		skillPath := path.Join(asset.DirAgySkills, name, asset.FileSKILLMd)
		body, err := FS.ReadFile(skillPath)
		if err != nil {
			t.Errorf("read %s: %v", skillPath, err)
			continue
		}

		lines := bytes.Split(body, []byte("\n"))
		for i, line := range lines {
			if bytes.HasPrefix(bytes.TrimSpace(line), []byte("allowed-tools:")) {
				t.Errorf(
					"%s:%d: carries allowed-tools: frontmatter",
					skillPath, i+1,
				)
			}
		}

		// Body after YAML frontmatter
		fmEnd := -1
		if len(lines) > 0 && bytes.Equal(bytes.TrimSpace(lines[0]), fence) {
			for i := 1; i < len(lines); i++ {
				if bytes.Equal(bytes.TrimSpace(lines[i]), fence) {
					fmEnd = i
					break
				}
			}
		}
		var contentAfterFM []byte
		if fmEnd != -1 && fmEnd+1 < len(lines) {
			contentAfterFM = bytes.Join(lines[fmEnd+1:], []byte("\n"))
		} else {
			contentAfterFM = body
		}

		// Antigravity strict validator limit: body <= 24,000 bytes
		if len(contentAfterFM) > 24000 {
			t.Errorf(
				"%s: body size %d exceeds Antigravity 24,000-byte limit",
				skillPath, len(contentAfterFM),
			)
		}

		// Check every cited references/... exists in embedded FS
		for _, ref := range reRef.FindAllString(string(body), -1) {
			refPath := path.Join(asset.DirAgySkills, name, ref)
			if _, statErr := fs.Stat(FS, refPath); statErr != nil {
				t.Errorf(
					"%s cites %s but it is not embedded: %v",
					skillPath, ref, statErr,
				)
			}
		}

		// If references/extended.md exists, verify its citations as well
		extendedPath := path.Join(
			asset.DirAgySkills, name, asset.DirReferences, "extended.md",
		)
		if extBody, extErr := FS.ReadFile(extendedPath); extErr == nil {
			for _, ref := range reRef.FindAllString(string(extBody), -1) {
				refPath := path.Join(asset.DirAgySkills, name, ref)
				if _, statErr := fs.Stat(FS, refPath); statErr != nil {
					t.Errorf(
						"%s cites %s but it is not embedded: %v",
						extendedPath, ref, statErr,
					)
				}
			}
		}
	}
}
