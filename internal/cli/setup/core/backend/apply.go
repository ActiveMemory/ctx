//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Setup is the entrypoint called by `ctx setup --backend
// <name>`. Resolves per-vendor defaults, merges user
// overrides, runs [Apply], and emits the user-facing
// confirmation through `write/setup`.
//
// Parameters:
//   - cmd: cobra command for stdout / stderr.
//   - name: backend type label (e.g. "vllm").
//   - endpoint: --endpoint override; empty uses per-vendor
//     default.
//   - apiKeyEnv: --api-key-env override; empty uses
//     per-vendor default.
//
// Returns:
//   - error: typed err/setup or err/context sentinel on
//     project resolution, .ctxrc IO, YAML, or write
//     failure.
func Setup(cmd *cobra.Command, name, endpoint, apiKeyEnv string) error {
	if name == "" {
		return errSetup.ErrBackendNameRequired
	}
	cfg := Resolve(name, endpoint, apiKeyEnv)
	ctxDir, dirErr := rc.ContextDir()
	if dirErr != nil {
		return dirErr
	}
	res, applyErr := Apply(filepath.Dir(ctxDir), cfg)
	if applyErr != nil {
		return applyErr
	}
	writeSetup.InfoBackendApplied(
		cmd, cfg.Name, res.Path, res.Created, res.Updated,
	)
	return nil
}

// Resolve merges per-vendor defaults with user
// --endpoint / --api-key-env overrides. Unknown backend
// names pass through with just the user-supplied fields
// so a custom name still works.
//
// Parameters:
//   - name: the backend type label.
//   - endpoint: user override or "".
//   - apiKeyEnv: user override or "".
//
// Returns:
//   - entity.BackendConfig: merged config ready for
//     [Apply].
func Resolve(name, endpoint, apiKeyEnv string) entity.BackendConfig {
	cfg := entity.BackendConfig{Name: name}
	switch name {
	case cfgBackend.NameVLLM:
		cfg.Endpoint = cfgBackend.DefaultEndpointVLLM
	case cfgBackend.NameOpenAI:
		cfg.Endpoint = cfgBackend.DefaultEndpointOpenAI
		cfg.APIKeyEnv = cfgBackend.EnvOpenAIAPIKey
	case cfgBackend.NameAnthropic:
		cfg.Endpoint = cfgBackend.DefaultEndpointAnthropic
		cfg.APIKeyEnv = cfgBackend.EnvAnthropicAPIKey
	case cfgBackend.NameOllama:
		cfg.Endpoint = cfgBackend.DefaultEndpointOllama
	case cfgBackend.NameLMStudio:
		cfg.Endpoint = cfgBackend.DefaultEndpointLMStudio
	}
	if endpoint != "" {
		cfg.Endpoint = endpoint
	}
	if apiKeyEnv != "" {
		cfg.APIKeyEnv = apiKeyEnv
	}
	return cfg
}

// Apply adds or updates a single entry in the `backends:`
// list of the project's `.ctxrc`. Idempotent: re-running
// with the same input is a no-op on disk if the entry
// is already byte-identical.
//
// Parameters:
//   - projectRoot: absolute path of the project root
//     (parent of `.context/`); `.ctxrc` lives at
//     `<projectRoot>/.ctxrc`.
//   - entry: the backend config to add or update,
//     addressed by [entity.BackendConfig.Name].
//
// Returns:
//   - Result: outcome (path, created, updated).
//   - error: typed err/setup sentinel on read / parse /
//     marshal / write failure.
func Apply(projectRoot string, entry entity.BackendConfig) (Result, error) {
	rcPath := filepath.Join(projectRoot, file.CtxRC)
	data, readErr := ctxIo.SafeReadUserFile(rcPath)
	created := false
	if readErr != nil {
		if !errors.Is(readErr, os.ErrNotExist) {
			return Result{}, errSetup.ReadCtxrc(rcPath, readErr)
		}
		created = true
		data = nil
	}
	doc, parseErr := parseDocument(data)
	if parseErr != nil {
		return Result{}, errSetup.ParseCtxrc(rcPath, parseErr)
	}
	updated, mutErr := upsert(doc, entry)
	if mutErr != nil {
		return Result{}, errSetup.MarshalCtxrc(mutErr)
	}
	out, marshalErr := yaml.Marshal(doc)
	if marshalErr != nil {
		return Result{}, errSetup.MarshalCtxrc(marshalErr)
	}
	if wErr := ctxIo.SafeWriteFile(rcPath, out, fs.PermFile); wErr != nil {
		return Result{}, errSetup.WriteFile(rcPath, wErr)
	}
	return Result{Path: rcPath, Created: created, Updated: updated}, nil
}
