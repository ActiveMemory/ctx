# Context CLI Makefile
#
# Common targets for Go developers

.PHONY: build ctxctl test vet fmt fmt-context lint lint-style lint-drift lint-shellcheck lint-powershell \
clean all release build-all help \
test-coverage smoke site site-guard site-feed site-serve site-serve-lan site-setup audit check plugin-reload \
journal journal-serve journal-serve-lan gpg-fix gpg-test register-mcp reinstall check-tools \
sync-version check-version-sync sync-why check-why sync-copilot-skills check-copilot-skills sync-codex-skills check-codex-skills codex-plugin-install sync-steering check-steering gemini-search \
gitnexus-version gitnexus-update gitnexus-index gitnexus-mcp strip-gitnexus install-ctxctl reinstall-ctxctl

# Default binary name and output
BINARY := ctx
OUTPUT := $(BINARY)

# Maintainer-only binary (separate Go module at tools/ctxctl,
# resolved via the repo-root go.work). Never shipped to users.
# Built into dist/ and installed to PATH alongside ctx, so every
# repo copy / worktree shares one binary and the root stays clean.
CTXCTL_BINARY := ctxctl
CTXCTL_OUTPUT := dist/$(CTXCTL_BINARY)

# Exact zensical pin for site builds, read from the tooling manifest
# (leading '=' stripped) so the installer and the check can never
# disagree about the version.
ZENSICAL_PIN := $(shell awk '$$1 == "zensical" && $$2 == "bin" { sub(/^=/, "", $$4); print $$4 }' hack/tool-versions.txt)

# Default target
all: build

## sync-version: Stamp VERSION into embedded plugin.json
sync-version:
	@V=$$(cat VERSION | tr -d '[:space:]'); \
	jq --arg v "$$V" '.version = $$v' internal/assets/claude/.claude-plugin/plugin.json > internal/assets/claude/.claude-plugin/plugin.json.tmp && \
	mv internal/assets/claude/.claude-plugin/plugin.json.tmp internal/assets/claude/.claude-plugin/plugin.json; \
	jq --arg v "$$V" '.version = $$v' internal/assets/codex/.codex-plugin/plugin.json > internal/assets/codex/.codex-plugin/plugin.json.tmp && \
	mv internal/assets/codex/.codex-plugin/plugin.json.tmp internal/assets/codex/.codex-plugin/plugin.json; \
	jq --arg v "$$V" '.version = $$v' internal/assets/claude/.codex-plugin/plugin.json > internal/assets/claude/.codex-plugin/plugin.json.tmp && \
	mv internal/assets/claude/.codex-plugin/plugin.json.tmp internal/assets/claude/.codex-plugin/plugin.json; \
	jq --arg v "$$V" '.metadata.version = $$v' .agents/plugins/marketplace.json > .agents/plugins/marketplace.json.tmp && \
	mv .agents/plugins/marketplace.json.tmp .agents/plugins/marketplace.json; \
	echo "Plugin version synced to $$V"

## build: Build for current platform (syncs version + embedded docs + copilot skills first)
build: sync-version sync-why sync-copilot-skills sync-codex-skills
	CGO_ENABLED=0 go build -ldflags="-X github.com/ActiveMemory/ctx/internal/bootstrap.version=$$(cat VERSION | tr -d '[:space:]')" -o $(OUTPUT) ./cmd/ctx

## ctxctl: Build the maintainer-only ctxctl binary (audit channel) into dist/
ctxctl:
	@mkdir -p dist
	CGO_ENABLED=0 go build -o $(CTXCTL_OUTPUT) ./tools/ctxctl

## test: Run tests with coverage summary
test:
	@CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test -cover ./...

## test-v: Run tests with verbose output
test-v:
	CGO_ENABLED=0 go test -v ./...

## test-cover: Generate HTML coverage report in dist/coverage.html
test-cover:
	@mkdir -p dist
	@CGO_ENABLED=0 go test -coverprofile=dist/coverage.out ./...
	@go tool cover -html=dist/coverage.out -o dist/coverage.html
	@echo "Coverage report: dist/coverage.html"

## test-coverage: Run tests with coverage and check against target (70%)
test-coverage:
	@echo "Running coverage check (target: 70%)..."
	@echo ""
	@CGO_ENABLED=0 go test -cover ./internal/context ./internal/cli 2>&1 | tee /tmp/ctx-coverage.txt
	@echo ""
	@CONTEXT_COV=$$(grep 'internal/context' /tmp/ctx-coverage.txt | grep -oE '[0-9]+\.[0-9]+%' | sed 's/%//'); \
	CLI_COV=$$(grep 'internal/cli' /tmp/ctx-coverage.txt | grep -oE '[0-9]+\.[0-9]+%' | sed 's/%//'); \
	echo "Coverage summary:"; \
	echo "  internal/context: $${CONTEXT_COV}% (target: 70%)"; \
	echo "  internal/cli: $${CLI_COV}% (target: 70% - aspirational)"; \
	echo ""; \
	if [ $$(echo "$$CONTEXT_COV < 70" | bc -l) -eq 1 ]; then \
		echo "FAIL: internal/context coverage below 70%"; \
		rm -f /tmp/ctx-coverage.txt; \
		exit 1; \
	fi; \
	echo "Coverage check passed (internal/context >= 70%)"; \
	rm -f /tmp/ctx-coverage.txt

## smoke: Build and run basic commands to verify binary works
smoke: build
	@echo "Running smoke tests..."
	@TMPDIR=$$(mktemp -d) && \
	cd $$TMPDIR && \
	echo "  Testing: ctx --help" && \
	$(CURDIR)/$(BINARY) --help > /dev/null && \
	echo "  Testing: ctx init" && \
	CTX_SKIP_PATH_CHECK=1 $(CURDIR)/$(BINARY) init > /dev/null && \
	echo "  Testing: ctx status" && \
	$(CURDIR)/$(BINARY) status > /dev/null && \
	echo "  Testing: ctx agent" && \
	$(CURDIR)/$(BINARY) agent > /dev/null && \
	echo "  Testing: ctx drift" && \
	$(CURDIR)/$(BINARY) drift > /dev/null && \
	echo "  Testing: ctx add task 'smoke test task'" && \
	$(CURDIR)/$(BINARY) add task "smoke test task" > /dev/null && \
	echo "  Testing: ctx journal source" && \
	$(CURDIR)/$(BINARY) journal source > /dev/null && \
	echo "  Testing: ctx why manifesto" && \
	$(CURDIR)/$(BINARY) why manifesto > /dev/null && \
	rm -rf $$TMPDIR && \
	echo "" && \
	echo "Smoke tests passed!"

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format code
fmt:
	go fmt ./...

## fmt-context: Format context files to 80-char line width
fmt-context:
	ctx fmt

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run

## lint-style: Run all cosmetic/style lint scripts (advisory, not fatal)
lint-style:
	@echo "==> Checking code drift..."
	@./hack/lint-drift.sh
	@echo "==> Checking docstrings..."
	@./hack/lint-docstrings.sh
	@echo "==> Checking mixed funcs..."
	@./hack/lint-mixed-funcs.sh
	@echo "==> Checking import conventions..."
	@./hack/lint-imports.sh
	@echo ""
	@echo "Style checks passed!"

## lint-drift: Check for code-level drift (magic strings, literal \n, Printf)
lint-drift:
	@./hack/lint-drift.sh

## lint-shellcheck: Run shellcheck on embedded *.sh scripts (warning+)
lint-shellcheck:
	@./hack/lint-shellcheck.sh

## lint-powershell: Run PSScriptAnalyzer on embedded *.ps1 scripts (Warning+)
lint-powershell:
	@./hack/lint-powershell.sh

## audit: Run all CI checks locally (fmt, vet, lint, drift, docs, test)
audit:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
	@echo "==> Running go vet..."
	@CGO_ENABLED=0 go vet ./...
	@echo "==> Running golangci-lint..."
	@golangci-lint run --timeout=5m
	@echo "==> Running style checks..."
	@$(MAKE) --no-print-directory lint-style
	@if command -v shellcheck >/dev/null 2>&1; then \
		echo "==> Running shellcheck..."; \
		$(MAKE) --no-print-directory lint-shellcheck; \
	else \
		echo "==> Skipping shellcheck (not installed locally; CI enforces this)"; \
	fi
	@if command -v pwsh >/dev/null 2>&1; then \
		echo "==> Running PSScriptAnalyzer..."; \
		$(MAKE) --no-print-directory lint-powershell; \
	else \
		echo "==> Skipping PSScriptAnalyzer (pwsh not installed locally; CI enforces this)"; \
	fi
	@echo "==> Checking version sync..."
	@$(MAKE) --no-print-directory check-version-sync
	@echo "==> Checking why docs freshness..."
	@$(MAKE) --no-print-directory check-why
	@echo "==> Checking Copilot skills freshness..."
	@$(MAKE) --no-print-directory check-copilot-skills
	@echo "==> Checking Codex skills freshness..."
	@$(MAKE) --no-print-directory check-codex-skills
	@echo "==> Checking steering outputs freshness..."
	@$(MAKE) --no-print-directory check-steering
	@echo "==> Running tests..."
	@CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test ./...
	@echo ""
	@echo "All checks passed!"
	@echo "Tip: run /ctx-link-check to verify doc links before committing."

## check: Build + audit (single entry point for build, fmt, vet, lint, test)
check: build audit

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -f $(CTXCTL_BINARY)
	rm -f tools/ctxctl/$(CTXCTL_BINARY)
	rm -rf dist/

## release: Full release process (build, tag, push)
release:
	./hack/release.sh

## build-all: Build binaries for all platforms (no tag)
build-all:
	./hack/build-all.sh $$(cat VERSION | tr -d '[:space:]')

## release-notes: Generate release notes (use Claude Code slash command)
release-notes:
	@echo "To generate release notes, run in Claude Code:"
	@echo ""
	@echo "  /release-notes"
	@echo ""
	@echo "This will analyze commits since the last tag and write to dist/RELEASE_NOTES.md"

## install: Install to /usr/local/bin (run as: make build && sudo make install)
install:
	@test -f $(BINARY) || (echo "Binary not found. Run 'make build' first, then 'sudo make install'" && exit 1)
	install -m 0755 $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed ctx to /usr/local/bin/ctx"

## reinstall: Build and install in one step
reinstall: build
	install -m 0755 $(BINARY) /usr/local/bin/$(BINARY) 2>/dev/null || sudo install -m 0755 $(BINARY) /usr/local/bin/$(BINARY)
	@echo "ctx reinstalled to /usr/local/bin/ctx"

## install-ctxctl: Install the maintainer-only ctxctl binary to /usr/local/bin
install-ctxctl:
	@test -f $(CTXCTL_OUTPUT) || (echo "Binary not found. Run 'make ctxctl' first, then 'make install-ctxctl'" && exit 1)
	install -m 0755 $(CTXCTL_OUTPUT) /usr/local/bin/$(CTXCTL_BINARY) 2>/dev/null || sudo install -m 0755 $(CTXCTL_OUTPUT) /usr/local/bin/$(CTXCTL_BINARY)
	@echo "Installed ctxctl to /usr/local/bin/$(CTXCTL_BINARY)"

## reinstall-ctxctl: Build and install ctxctl in one step (maintainer-only)
reinstall-ctxctl: ctxctl
	install -m 0755 $(CTXCTL_OUTPUT) /usr/local/bin/$(CTXCTL_BINARY) 2>/dev/null || sudo install -m 0755 $(CTXCTL_OUTPUT) /usr/local/bin/$(CTXCTL_BINARY)
	@echo "ctxctl reinstalled to /usr/local/bin/$(CTXCTL_BINARY)"

## site-setup: Install the pinned zensical via pipx
site-setup:
	pipx install 'zensical==$(ZENSICAL_PIN)'

# Refuse to build the site with a generator that is not the pinned one.
# Both directions matter: behind the pin silently regenerates pages from
# an older generator, ahead of it regenerates them from a newer one. The
# size of the resulting diff varies with the version gap — it is not
# predictable, which is exactly why it must not be produced by accident.
# Reuses the manifest pin and the version comparison in check-tools.sh
# rather than reimplementing either.
site-guard:
	@./hack/check-tools.sh --only zensical --strict || { \
	  echo ""; \
	  echo "make site refused to run: zensical is not the pinned $(ZENSICAL_PIN)."; \
	  echo ""; \
	  echo "Rebuilding with a different generator rewrites pages in site/ with"; \
	  echo "that generator's output. CI never rebuilds the site, so this check"; \
	  echo "is the only thing standing between an accidental churn diff and main."; \
	  echo ""; \
	  echo "  Behind the pin  -> install the pin:  make site-setup"; \
	  echo "  Ahead of it     -> intentional upgrade? bump the pin in"; \
	  echo "                     hack/tool-versions.txt, rebuild, and commit the"; \
	  echo "                     regenerated site/ with the pin bump."; \
	  exit 1; }

## site: Build documentation site and generate feed
site: site-guard
	zensical build
	ctx site feed

## site-feed: Generate Atom feed from blog posts
site-feed:
	ctx site feed

## site-serve: Serve documentation site locally
site-serve:
	zensical serve

## site-serve-lan: Serve docs site on all interfaces (LAN-accessible)
site-serve-lan:
	zensical serve -a 0.0.0.0:8000

## journal: Import sessions and regenerate journal site
journal:
	@echo "==> Importing sessions to journal..."
	@ctx journal import --all
	@echo "==> Generating journal site..."
	@ctx journal site --build
	@echo ""
	@echo "Journal site updated!"
	@echo ""
	@echo "Next steps (in Claude Code):"
	@echo "  /ctx-journal-enrich-all  — exports if needed + adds metadata per entry"
	@echo ""
	@echo "Then re-run: make journal"

## journal-serve: Serve the journal site (port 8001; docs uses 8000)
journal-serve:
	@ctx journal site
	cd .context/journal-site && zensical serve -a localhost:8001

## journal-serve-lan: Serve journal site on all interfaces (LAN-accessible, port 8001)
journal-serve-lan:
	cd .context/journal-site && zensical serve -a 0.0.0.0:8001

## gpg-fix: Fix GPG signing configuration
gpg-fix:
	./hack/gpg-fix.sh

## gpg-test: Test GPG signing configuration
gpg-test:
	./hack/gpg-fix.sh --test

## check-tools: Verify tooling dependency versions (Go, Node, GitNexus, MCP servers, ...)
check-tools:
	@./hack/check-tools.sh

## register-mcp: Register all MCP servers (gemini-search, gitnexus) with Claude Code
register-mcp:
	@./hack/register-gemini-search.sh
	@./hack/register-gitnexus.sh

## gitnexus-version: Check for gitnexus version drift
gitnexus-version:
	@INSTALLED=$$(gitnexus --version 2>/dev/null || echo "not installed"); \
	LATEST=$$(npm view gitnexus version 2>/dev/null || echo "unknown"); \
	echo "Installed: $$INSTALLED"; \
	echo "Latest:    $$LATEST"; \
	if [ "$$INSTALLED" = "$$LATEST" ]; then \
		echo "Up to date."; \
	else \
		echo "Update available — run 'make gitnexus-update'"; \
	fi

## gitnexus-update: Update gitnexus to latest version
gitnexus-update:
	npm install -g gitnexus@latest
	@echo "Updated to $$(gitnexus --version)"

## gitnexus-analyze: Updates gitnexus embeddings and skill.
gitnexus-analyze:
	gitnexus analyze --embeddings --skill
	echo "GitNexus updated AGENTS.md and CLAUDE.md -- DO NOT COMMIT THEM!"

## gitnexus-index: Index this repo into GitNexus (Docker; no npm binary needed)
gitnexus-index:
	@./hack/gitnexus-index.sh

## gitnexus-mcp: Launch the GitNexus stdio MCP server (what `claude mcp add` registers)
gitnexus-mcp:
	@./hack/gitnexus-docker.sh mcp

## strip-gitnexus: Remove the GitNexus auto-injected block from AGENTS.md/CLAUDE.md
strip-gitnexus:
	@./hack/strip-gitnexus.sh

## gemini-search: Register gemini-search MCP server with Claude Code
gemini-search:
	@./hack/register-gemini-search.sh

## plugin-reload: Clear cached plugin (restart Claude Code to pick up skill/hook changes)
plugin-reload:
	@./hack/plugin-reload.sh

## sync-why: Copy philosophy docs into internal/assets/why/ for embedding
sync-why:
	cp docs/index.md internal/assets/why/manifesto.md
	cp docs/home/about.md internal/assets/why/about.md
	cp docs/reference/design-invariants.md internal/assets/why/design-invariants.md
	@echo "Why docs synced."

## check-version-sync: Verify VERSION file matches embedded plugin.json
check-version-sync:
	@V=$$(cat VERSION | tr -d '[:space:]'); \
	PV=$$(jq -r '.version' internal/assets/claude/.claude-plugin/plugin.json); \
	if [ "$$V" != "$$PV" ]; then \
		echo "FAIL: VERSION ($$V) != plugin.json ($$PV) — run 'make sync-version'"; \
		exit 1; \
	fi; \
	CV=$$(jq -r '.version' internal/assets/codex/.codex-plugin/plugin.json); \
	if [ "$$V" != "$$CV" ]; then \
		echo "FAIL: VERSION ($$V) != codex plugin.json ($$CV) — run 'make sync-version'"; \
		exit 1; \
	fi; \
	DV=$$(jq -r '.version' internal/assets/claude/.codex-plugin/plugin.json); \
	if [ "$$V" != "$$DV" ]; then \
		echo "FAIL: VERSION ($$V) != claude-root codex plugin.json ($$DV) — run 'make sync-version'"; \
		exit 1; \
	fi; \
	MV=$$(jq -r '.metadata.version' .agents/plugins/marketplace.json); \
	if [ "$$V" != "$$MV" ]; then \
		echo "FAIL: VERSION ($$V) != .agents/plugins/marketplace.json ($$MV) — run 'make sync-version'"; \
		exit 1; \
	fi; \
	echo "Version sync OK ($$V)."

## sync-copilot-skills: Sync Copilot CLI skills from canonical ctx skills
sync-copilot-skills:
	@./hack/sync-copilot-skills.sh

## sync-steering: Regenerate tool-native steering outputs from .context/steering
sync-steering:
	@CGO_ENABLED=0 go run ./cmd/ctx steering sync --all

## check-steering: Verify tracked steering outputs match .context/steering source
check-steering:
	@CGO_ENABLED=0 go run ./cmd/ctx steering sync --all > /dev/null
	@if ! git diff --quiet -- .cursor .clinerules .kiro/steering; then \
		echo "FAIL: steering outputs are stale — run 'make sync-steering' and commit"; \
		git --no-pager diff --stat -- .cursor .clinerules .kiro/steering; \
		exit 1; \
	fi
	@echo "Steering outputs are in sync."

## check-copilot-skills: Verify Copilot CLI skills match ctx source skills
check-copilot-skills:
	@TMPDIR=$$(mktemp -d) && \
	cp -r internal/assets/integrations/copilot-cli/skills/ "$$TMPDIR/before" && \
	./hack/sync-copilot-skills.sh > /dev/null && \
	if ! diff -rq "$$TMPDIR/before" internal/assets/integrations/copilot-cli/skills/ > /dev/null 2>&1; then \
		echo "FAIL: Copilot CLI skills are stale — run 'make sync-copilot-skills'"; \
		diff -rq "$$TMPDIR/before" internal/assets/integrations/copilot-cli/skills/ || true; \
		cp -r "$$TMPDIR/before/"* internal/assets/integrations/copilot-cli/skills/; \
		rm -rf "$$TMPDIR"; \
		exit 1; \
	fi; \
	rm -rf "$$TMPDIR"; \
	echo "Copilot CLI skills are in sync."

## sync-codex-skills: Sync Codex plugin skills from canonical ctx skills
sync-codex-skills:
	@./hack/sync-codex-skills.sh

## check-codex-skills: Verify Codex plugin skills match ctx source skills
check-codex-skills:
	@TMPDIR=$$(mktemp -d) && \
	cp -r internal/assets/codex/skills/ "$$TMPDIR/before" && \
	cp internal/assets/claude/hooks/codex.json "$$TMPDIR/codex.json" && \
	cp internal/assets/claude/.codex-plugin/plugin.json "$$TMPDIR/plugin.json" && \
	./hack/sync-codex-skills.sh > /dev/null && \
	if ! diff -rq "$$TMPDIR/before" internal/assets/codex/skills/ > /dev/null 2>&1 || \
	   ! diff -q "$$TMPDIR/codex.json" internal/assets/claude/hooks/codex.json > /dev/null 2>&1 || \
	   ! diff -q "$$TMPDIR/plugin.json" internal/assets/claude/.codex-plugin/plugin.json > /dev/null 2>&1; then \
		echo "FAIL: Codex skills or dual-manifest files are stale — run 'make sync-codex-skills'"; \
		diff -rq "$$TMPDIR/before" internal/assets/codex/skills/ || true; \
		diff -q "$$TMPDIR/codex.json" internal/assets/claude/hooks/codex.json || true; \
		diff -q "$$TMPDIR/plugin.json" internal/assets/claude/.codex-plugin/plugin.json || true; \
		rm -rf internal/assets/codex/skills && cp -r "$$TMPDIR/before" internal/assets/codex/skills; \
		cp "$$TMPDIR/codex.json" internal/assets/claude/hooks/codex.json; \
		cp "$$TMPDIR/plugin.json" internal/assets/claude/.codex-plugin/plugin.json; \
		rm -rf "$$TMPDIR"; \
		exit 1; \
	fi; \
	rm -rf "$$TMPDIR"; \
	echo "Codex skills are in sync."

## codex-plugin-install: Register this checkout as a Codex marketplace and install the ctx plugin
codex-plugin-install:
	@codex plugin marketplace add "$$(pwd)" && codex plugin add ctx@activememory-ctx
	@echo "Open codex and run /hooks to review and trust the ctx hooks."

## check-why: Verify embedded why docs match source docs
check-why:
	@diff -q docs/index.md internal/assets/why/manifesto.md || (echo "FAIL: manifesto.md is stale — run 'make sync-why'" && exit 1)
	@diff -q docs/home/about.md internal/assets/why/about.md || (echo "FAIL: about.md is stale — run 'make sync-why'" && exit 1)
	@diff -q docs/reference/design-invariants.md internal/assets/why/design-invariants.md || (echo "FAIL: design-invariants.md is stale — run 'make sync-why'" && exit 1)
	@echo "Why docs are in sync."

## title-case-check: Dry-run title-case checker on docs (or TARGET=path)
title-case-check:
	@python3 hack/title-case-headings.py $${TARGET:-docs}

## title-case-fix: Apply title-case fixes to headings + admonition titles (TARGET=path defaults to docs)
title-case-fix:
	@python3 hack/title-case-headings.py --apply $${TARGET:-docs}

## help: Show this help
help:
	@echo "Context CLI - Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

-include Makefile.ctx
