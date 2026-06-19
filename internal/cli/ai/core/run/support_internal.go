//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package run

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	backendPkg "github.com/ActiveMemory/ctx/internal/backend"
	cfgAI "github.com/ActiveMemory/ctx/internal/config/ai"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	cfgFS "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// resolve returns the selected configured backend.
//
// Parameters:
//   - backendName: optional backend selector
//
// Returns:
//   - resolvedBackend: selected backend and metadata
//   - error: rc or registry resolution error
func resolve(backendName string) (resolvedBackend, error) {
	registry := &backendPkg.Registry{}
	configs := rc.RC().Backends.Configs
	if len(configs) == 0 {
		_, resolveErr := registry.Default()
		return resolvedBackend{}, resolveErr
	}
	for name, config := range configs {
		backendConfig := backendPkg.Config{
			Name:         name,
			Type:         config.Type,
			Endpoint:     config.Endpoint,
			APIKeyEnv:    config.APIKeyEnv,
			Timeout:      config.Timeout,
			DefaultModel: config.DefaultModel,
		}
		registerErr := registry.RegisterBuiltin(name, backendConfig)
		if registerErr != nil {
			return resolvedBackend{}, registerErr
		}
	}
	registry.SetDefault(rc.RC().Backends.Default)
	selectedName := selectedBackendName(backendName, configs)
	selected, resolveErr := resolveSelected(registry, selectedName)
	if resolveErr != nil {
		return resolvedBackend{}, resolveErr
	}
	return resolvedBackend{
		name:    selectedName,
		backend: selected,
	}, nil
}

// resolveSelected resolves either an explicit or registry default backend.
//
// Parameters:
//   - registry: backend registry
//   - selectedName: explicit selected backend name
//
// Returns:
//   - backend.Backend: resolved backend
//   - error: registry resolution failure
func resolveSelected(
	registry *backendPkg.Registry,
	selectedName string,
) (backendPkg.Backend, error) {
	if selectedName != "" {
		return registry.Resolve(selectedName)
	}
	return registry.Default()
}

// selectedBackendName returns the configured backend key to resolve.
//
// Parameters:
//   - backendName: explicit backend selector
//   - configs: configured backend map
//
// Returns:
//   - string: backend key, or empty when registry default should decide
func selectedBackendName(
	backendName string,
	configs map[string]rc.BackendRC,
) string {
	if backendName != "" {
		return backendName
	}
	defaultName := rc.RC().Backends.Default
	if defaultName != "" {
		return defaultName
	}
	if len(configs) != 1 {
		return ""
	}
	for name := range configs {
		return name
	}
	return ""
}

// writeArtifact writes a proposed-patch artifact under .context.
//
// Parameters:
//   - artifact: proposal artifact to persist
//
// Returns:
//   - string: artifact path
//   - error: marshal, mkdir, or write failure
func writeArtifact(artifact ProposalArtifact) (string, error) {
	indent := token.Space + token.Space
	data, marshalErr := json.MarshalIndent(artifact, "", indent)
	if marshalErr != nil {
		return "", errBackend.BadRequest{
			Name:  artifact.Backend,
			Cause: marshalErr,
		}
	}
	proposalDir := filepath.Join(dir.Context, cfgAI.DirProposals, cfgAI.DirAI)
	mkdirErr := ctxio.SafeMkdirAll(proposalDir, cfgFS.PermExec)
	if mkdirErr != nil {
		return "", mkdirErr
	}
	name := cfgAI.ArtifactPrefix +
		time.Now().UTC().Format(cfgAI.TimestampLayout) +
		cfgAI.ArtifactExtJSON
	path := filepath.Join(proposalDir, name)
	writeErr := ctxio.SafeWriteFile(path, data, cfgFS.PermFile)
	if writeErr != nil {
		return "", writeErr
	}
	return path, nil
}

// splitEmit splits comma-separated emit kinds.
//
// Parameters:
//   - emit: comma-separated emit kinds
//
// Returns:
//   - []string: trimmed emit kinds
func splitEmit(emit string) []string {
	parts := strings.Split(emit, cfgAI.EmitSeparator)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
