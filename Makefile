# Context CLI Makefile
#
# Common targets for Go developers

.PHONY: build test vet fmt lint lint-drift lint-docs clean all release build-all dogfood help \
test-coverage smoke site site-feed site-serve site-serve-lan site-setup audit check plugin-reload \
journal journal-serve journal-serve-lan watch-session backup backup-global backup-all gpg-fix gpg-test \
sync-why check-why

# Default binary name and output
BINARY := ctx
OUTPUT := $(BINARY)

# Default target
all: build

## build: Build for current platform
build:
	CGO_ENABLED=0 go build -ldflags="-X github.com/ActiveMemory/ctx/internal/bootstrap.version=$$(cat VERSION | tr -d '[:space:]')" -o $(OUTPUT) ./cmd/ctx

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
	echo "  Testing: ctx recall list" && \
	$(CURDIR)/$(BINARY) recall list > /dev/null && \
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

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run

## lint-drift: Check for code-level drift (magic strings, literal \n, Printf)
lint-drift:
	@./hack/lint-drift.sh

## lint-docs: Check doc.go file listings match actual files
lint-docs:
	@./hack/lint-docs.sh

## audit: Run all CI checks locally (fmt, vet, lint, drift, docs, test)
audit:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
	@echo "==> Running go vet..."
	@CGO_ENABLED=0 go vet ./...
	@echo "==> Running golangci-lint..."
	@golangci-lint run --timeout=5m
	@echo "==> Checking code drift..."
	@./hack/lint-drift.sh
	@echo "==> Checking doc.go listings..."
	@./hack/lint-docs.sh
	@echo "==> Checking why docs freshness..."
	@$(MAKE) --no-print-directory check-why
	@echo "==> Running tests..."
	@CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test ./...
	@echo ""
	@echo "All checks passed!"

## check: Build + audit (single entry point for build, fmt, vet, lint, test)
check: build audit

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
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

## dogfood: Start dogfooding in a target folder
dogfood:
	@test -n "$(TARGET)" || (echo "Usage: make dogfood TARGET=~/WORKSPACE/ctx-dogfood" && exit 1)
	./hack/start-dogfood.sh $(TARGET)

## install: Install to /usr/local/bin (run as: make build && sudo make install)
install:
	@test -f $(BINARY) || (echo "Binary not found. Run 'make build' first, then 'sudo make install'" && exit 1)
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed ctx to /usr/local/bin/ctx"

## site-setup: Install zensical via pipx
site-setup:
	pipx install zensical

## site: Build documentation site and generate feed
site:
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

## journal: Export sessions and regenerate journal site
journal:
	@echo "==> Exporting sessions to journal..."
	@ctx recall export --all
	@echo "==> Generating journal site..."
	@ctx journal site --build
	@echo ""
	@echo "Journal site updated!"
	@echo ""
	@echo "Next steps (in Claude Code):"
	@echo "  /ctx-journal-enrich-all  — add metadata per entry (skips if frontmatter exists)"
	@echo ""
	@echo "Then re-run: make journal"

## journal-serve: Serve the journal site (port 8001; docs uses 8000)
journal-serve:
	@ctx journal site
	cd .context/journal-site && zensical serve -a localhost:8001

## journal-serve-lan: Serve journal site on all interfaces (LAN-accessible, port 8001)
journal-serve-lan:
	cd .context/journal-site && zensical serve -a 0.0.0.0:8001

## backup: Backup project context (.context/ and .claude/) to SMB share
backup:
	./hack/backup-context.sh

## backup-global: Backup global Claude Code data (~/.claude/) to SMB share
backup-global:
	./hack/backup-global.sh

## backup-all: Backup both project context and global Claude data
backup-all: backup backup-global

## gpg-fix: Fix GPG signing configuration
gpg-fix:
	./hack/gpg-fix.sh

## gpg-test: Test GPG signing configuration
gpg-test:
	./hack/gpg-fix.sh --test

## plugin-reload: Clear cached plugin (restart Claude Code to pick up skill/hook changes)
plugin-reload:
	@./hack/plugin-reload.sh

## sync-why: Copy philosophy docs into internal/assets/why/ for embedding
sync-why:
	cp docs/index.md internal/assets/why/manifesto.md
	cp docs/home/about.md internal/assets/why/about.md
	cp docs/reference/design-invariants.md internal/assets/why/design-invariants.md
	@echo "Why docs synced."

## check-why: Verify embedded why docs match source docs
check-why:
	@diff -q docs/index.md internal/assets/why/manifesto.md || (echo "FAIL: manifesto.md is stale — run 'make sync-why'" && exit 1)
	@diff -q docs/home/about.md internal/assets/why/about.md || (echo "FAIL: about.md is stale — run 'make sync-why'" && exit 1)
	@diff -q docs/reference/design-invariants.md internal/assets/why/design-invariants.md || (echo "FAIL: design-invariants.md is stale — run 'make sync-why'" && exit 1)
	@echo "Why docs are in sync."

## watch-session: Watch current session for token usage
watch-session:
	./hack/context-watch.sh

## help: Show this help
help:
	@echo "Context CLI - Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

-include Makefile.ctx
