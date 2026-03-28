.PHONY: help build build-macos-arm64 build-macos-amd64 build-windows-arm64 build-windows-amd64 build-linux-arm64 build-linux-amd64 build-all clean test acceptance release man man-install install uninstall wiki

# Default target
.DEFAULT_GOAL := build

# Detect local platform
LOCAL_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/darwin/;s/linux/linux/')
LOCAL_ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

# Variables
BINARY_NAME := proj
GO := go
GOOS ?= $(LOCAL_OS)
GOARCH ?= $(LOCAL_ARCH)
OUTPUT_DIR := ./bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Platform-specific extensions
ifeq ($(GOOS),windows)
BINARY_EXT := .exe
endif

# Help target
help:
	@echo "Makefile targets for $(BINARY_NAME):"
	@echo ""
	@echo "Build targets:"
	@echo "  make build                    Build for local platform ($(LOCAL_OS)/$(LOCAL_ARCH))"
	@echo "  make build-macos-arm64        Build for macOS/arm64"
	@echo "  make build-macos-amd64        Build for macOS/amd64 (Intel)"
	@echo "  make build-windows-arm64      Build for Windows/arm64"
	@echo "  make build-windows-amd64      Build for Windows/amd64 (Intel)"
	@echo "  make build-linux-arm64        Build for Linux/arm64"
	@echo "  make build-linux-amd64        Build for Linux/amd64 (Intel)"
	@echo "  make build-all                Build for all platforms"
	@echo ""
	@echo "Development targets:"
	@echo "  make test                     Run unit tests"
	@echo "  make acceptance               Run acceptance tests (builds + tests)"
	@echo "  make acceptance-no-cleanup    Run acceptance tests and preserve temp directories"
	@echo "  make clean                    Remove build artifacts and acceptance temp directories"
	@echo "  make man                      Generate man pages from markdown sources"
	@echo "  make wiki                     Generate wiki markdown files from doc sources"
	@echo "  make install                  Install binary and man pages to ~/.local"
	@echo "  make uninstall                Remove binary and man pages from ~/.local"
	@echo "  make help                     Display this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  GOOS=<os>                     Target OS (darwin, windows, linux)"
	@echo "  GOARCH=<arch>                 Target architecture (arm64, amd64)"
	@echo ""
	@echo "Release targets:"
	@echo "  make release VERSION [COMMIT] [DRY]  Create GitHub PR-based release (full workflow)"
	@echo ""
	@echo "Examples:"
	@echo "  make build-windows-amd64             Build Windows 64-bit"
	@echo "  make build-all                       Build all platform variants"
	@echo "  make release v1.0.0                  Release from main branch"
	@echo "  make release v1.0.0 abc123de         Release from specific commit"
	@echo "  make release v1.0.0 DRY              Preview what would happen"
	@echo "  make release v1.0.0 abc123de DRY     Preview from specific commit"
	@echo "  GOOS=linux GOARCH=arm64 make         Custom build"

# Main build target (builds for local platform + creates copy as 'proj')
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT) .
	@cp $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT) $(OUTPUT_DIR)/$(BINARY_NAME)
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT)"
	@echo "✓ Copied to: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Platform-specific build targets
build-macos-arm64:
	@echo "Building $(BINARY_NAME) for macOS/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64"

build-macos-amd64:
	@echo "Building $(BINARY_NAME) for macOS/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64"

build-windows-arm64:
	@echo "Building $(BINARY_NAME) for Windows/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe"

build-windows-amd64:
	@echo "Building $(BINARY_NAME) for Windows/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-linux-arm64:
	@echo "Building $(BINARY_NAME) for Linux/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64"

build-linux-amd64:
	@echo "Building $(BINARY_NAME) for Linux/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "✓ Built: $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64"

# Build all platforms
build-all: build-macos-arm64 build-macos-amd64 build-windows-arm64 build-windows-amd64 build-linux-arm64 build-linux-amd64
	@echo ""
	@echo "✓ Successfully built for all platforms"
	@echo "Binaries in: $(OUTPUT_DIR)/"
	@ls -lh $(OUTPUT_DIR)/

# Test target
test: build
	@echo "Running tests..."
	$(GO) test -v ./...

# Acceptance test target
acceptance: build
	@echo "Running acceptance tests..."
	$(GO) test -v -tags=acceptance ./tests/acceptance/...

# Acceptance test target without cleanup (useful for debugging)
acceptance-no-cleanup: build
	@echo "Running acceptance tests without cleanup..."
	NO_CLEANUP=1 $(GO) test -v -tags=acceptance ./tests/acceptance/...

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleaning acceptance test temp directories..."
	@find /tmp -name "proj-acceptance-*" -type d -exec rm -rf {} + 2>/dev/null || true
	@find /var/folders -name "proj-acceptance-*" -type d -exec rm -rf {} + 2>/dev/null || true
	@$(GO) clean
	@echo "✓ Clean complete"

# Man page generation
MAN_SRC_DIR := ./doc
MAN_OUT_DIR := ./bin/man

# Wiki generation
WIKI_SRC_DIR := ./doc
WIKI_OUT_DIR := ./bin/wiki

man:
	@echo "Generating man pages..."
	@mkdir -p $(MAN_OUT_DIR)
	@command -v go-md2man >/dev/null 2>&1 || go install github.com/cpuguy83/go-md2man/v2@latest
	@for src in $(MAN_SRC_DIR)/*.md; do \
		out=$$(basename "$$src" .md); \
		go-md2man -in "$$src" -out "$(MAN_OUT_DIR)/$$out"; \
		echo "✓ Generated: $(MAN_OUT_DIR)/$$out"; \
	done
	@echo "✓ Man pages generated in $(MAN_OUT_DIR)/"

man-install: man
	@echo "Installing man pages..."
	@mkdir -p ~/.local/share/man/man1 ~/.local/share/man/man5 ~/.local/share/man/man7
	@if [ -f $(MAN_OUT_DIR)/proj.1 ]; then \
		cp $(MAN_OUT_DIR)/proj.1 ~/.local/share/man/man1/; \
		echo "✓ Installed: proj.1"; \
	fi
	@if [ -f $(MAN_OUT_DIR)/proj.yml.5 ]; then \
		cp $(MAN_OUT_DIR)/proj.yml.5 ~/.local/share/man/man5/; \
		echo "✓ Installed: proj.yml.5"; \
	fi
	@if [ -f $(MAN_OUT_DIR)/proj.7 ]; then \
		cp $(MAN_OUT_DIR)/proj.7 ~/.local/share/man/man7/; \
		echo "✓ Installed: proj.7"; \
	fi
	@echo "Updating man database..."
	@mandb ~/.local/share/man 2>/dev/null || true
	@echo "✓ Man pages installed. Try: man proj"

wiki:
	@echo "Generating wiki files..."
	@mkdir -p $(WIKI_OUT_DIR)
	@for src in $(WIKI_SRC_DIR)/*.md; do \
		base=$$(basename "$$src"); \
		if [ "$$base" = "proj.1.md" ]; then \
			dest="proj-command.md"; \
		elif [ "$$base" = "proj.7.md" ]; then \
			dest="proj-file-format.md"; \
		else \
			dest=$$(echo "$$base" | sed -E 's/\.[157]\.md$$/.md/'); \
		fi; \
		cp "$$src" "$(WIKI_OUT_DIR)/$$dest"; \
		echo "✓ Copied $$base to $(WIKI_OUT_DIR)/$$dest"; \
	done
	@echo "✓ Wiki documentation generated in $(WIKI_OUT_DIR)/"

# Install binary and man pages (XDG-compliant, user-level)
install: build man
	@echo "Installing proj..."
	@mkdir -p ~/.local/bin
	@cp $(OUTPUT_DIR)/$(BINARY_NAME) ~/.local/bin/
	@chmod +x ~/.local/bin/$(BINARY_NAME)
	@echo "✓ Installed binary to ~/.local/bin/$(BINARY_NAME)"
	@mkdir -p ~/.local/share/man/man1 ~/.local/share/man/man5 ~/.local/share/man/man7
	@cp $(MAN_OUT_DIR)/*.1 ~/.local/share/man/man1/
	@cp $(MAN_OUT_DIR)/*.5 ~/.local/share/man/man5/
	@cp $(MAN_OUT_DIR)/*.7 ~/.local/share/man/man7/
	@echo "✓ Installed man pages"
	@mandb ~/.local/share/man 2>/dev/null || true
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH."
	@echo "Add to your shell config (e.g., ~/.zshrc or ~/.bashrc):"
	@echo '  export PATH="$$HOME/.local/bin:$$PATH"'
	@echo ""
	@echo "Try: proj --version"

# Uninstall binary and man pages
uninstall:
	@echo "Uninstalling proj..."
	@rm -f ~/.local/bin/$(BINARY_NAME)
	@echo "✓ Removed binary from ~/.local/bin/"
	@rm -f ~/.local/share/man/man1/proj.1
	@rm -f ~/.local/share/man/man1/proj-new.1
	@rm -f ~/.local/share/man/man1/proj-add.1
	@rm -f ~/.local/share/man/man1/proj-install.1
	@rm -f ~/.local/share/man/man1/proj-info.1
	@rm -f ~/.local/share/man/man1/proj-help.1
	@rm -f ~/.local/share/man/man5/proj.yml.5
	@rm -f ~/.local/share/man/man7/proj.7
	@rm -f ~/.local/share/man/man7/proj-scripts.7
	@rm -f ~/.local/share/man/man7/proj-lua.7
	@rm -f ~/.local/share/man/man7/proj-variables.7
	@rm -f ~/.local/share/man/man7/proj-template.7
	@echo "✓ Removed man pages"
	@mandb ~/.local/share/man 2>/dev/null || true
	@echo "✓ Uninstall complete"

# Unified GitHub PR-based release target
release:
	@bash -c 'VERSION="$(word 2,$(MAKECMDGOALS))"; \
	COMMIT="$(word 3,$(MAKECMDGOALS))"; \
	DRY_RUN="$(word 4,$(MAKECMDGOALS))"; \
	if [ -z "$$VERSION" ] || [ "$$VERSION" = "release" ]; then \
		echo "ERROR: No version provided"; \
		echo "Usage: make release VERSION [COMMIT] [DRY]"; \
		echo "Examples:"; \
		echo "  make release v1.0.0"; \
		echo "  make release v1.0.0 abc123de"; \
		echo "  make release v1.0.0 DRY"; \
		echo "  make release v1.0.0 abc123de DRY"; \
		exit 1; \
	fi; \
	if [ -n "$$COMMIT" ] && [ "$$COMMIT" = "DRY" ]; then \
		DRY_RUN="$$COMMIT"; \
		COMMIT=""; \
	fi; \
	if [ -z "$$COMMIT" ]; then \
		COMMIT="main"; \
	fi; \
	if ! git rev-parse -q --verify "$$COMMIT" >/dev/null 2>&1; then \
		echo "ERROR: Invalid commit: $$COMMIT"; \
		exit 1; \
	fi; \
	if ! echo "$$VERSION" | grep -qE "^v[0-9]+\.[0-9]+\.[0-9]+$$"; then \
		echo "ERROR: Invalid version format: $$VERSION"; \
		echo "Must be semantic versioning format: vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi; \
	if git rev-parse -q --verify "refs/tags/$$VERSION" >/dev/null 2>&1; then \
		echo "ERROR: Tag $$VERSION already exists"; \
		echo "To see existing tags, run: git tag -l"; \
		exit 1; \
	fi; \
	if ! command -v gh &> /dev/null; then \
		echo "ERROR: GitHub CLI not found"; \
		echo "Install from: https://cli.github.com"; \
		exit 1; \
	fi; \
	if [ "$$DRY_RUN" = "DRY" ]; then \
		echo "DRY-RUN MODE (no changes will be made)"; \
		echo ""; \
		echo "Release: $$VERSION"; \
		echo "From commit: $$COMMIT"; \
		echo ""; \
		echo "Would perform:"; \
		echo "  1. Create branch release/$$VERSION (on $$COMMIT)"; \
		echo "  2. Update first line of README.md to: # proj $$VERSION"; \
		echo "  3. Commit change with message: tagging new release"; \
		echo "  4. Push branch to GitHub"; \
		echo "  5. Create PR with title: RELEASE: $$VERSION"; \
		echo "  6. Wait for CI checks (2 minute timeout)"; \
		echo "  7. Merge PR (squash)"; \
		echo "  8. Wait for merge to complete"; \
		echo "  9. Tag created via GitHub Actions"; \
		echo "  10. Delete local branch release/$$VERSION"; \
		echo ""; \
		echo "To execute, run: make release $$VERSION"; \
	else \
		echo "Creating release $$VERSION via GitHub PR..."; \
		echo "From commit: $$COMMIT"; \
		echo ""; \
		echo "Step 1/10: Creating branch..."; \
		BRANCH="release/$$VERSION"; \
		git checkout -b "$$BRANCH" "$$COMMIT" >/dev/null 2>&1 || { echo "ERROR: Failed to create branch on $$COMMIT"; exit 1; }; \
		echo "✓ Created branch $$BRANCH from $$COMMIT"; \
		echo ""; \
		echo "Step 2/10: Updating README.md..."; \
		if [ ! -f "README.md" ]; then echo "ERROR: README.md not found"; exit 1; fi; \
		sed -i.bak "1s/^#.*[Pp]roj.*$$/# Proj $$VERSION/" README.md || { echo "ERROR: Failed to update README.md"; exit 1; }; \
		rm -f README.md.bak; \
		echo "✓ Updated README.md first line"; \
		echo ""; \
		echo "Step 3/10: Committing change..."; \
		git add README.md || { echo "ERROR: Failed to add README.md"; exit 1; }; \
		git commit -m "tagging new release" >/dev/null 2>&1 || { echo "ERROR: Failed to commit"; exit 1; }; \
		echo "✓ Committed change"; \
		echo ""; \
		echo "Step 4/10: Pushing to GitHub..."; \
		git push -u origin "$$BRANCH" >/dev/null 2>&1 || { echo "ERROR: Failed to push branch"; exit 1; }; \
		echo "✓ Pushed to origin"; \
		echo ""; \
		echo "Step 5/10: Creating pull request..."; \
		PR_URL=$$(gh pr create --title "RELEASE: $$VERSION" --body "Release $$VERSION" --base main --head "$$BRANCH" 2>&1) || { echo "ERROR: Failed to create PR"; exit 1; }; \
		PR_NUMBER=$$(echo "$$PR_URL" | awk -F"/" "{print \$$NF}"); \
		echo "✓ Created PR #$$PR_NUMBER"; \
		echo "  Link: $$PR_URL"; \
		echo ""; \
		echo "Step 6/10: Waiting for CI checks (2 min timeout)..."; \
		CHECK_COUNT=0; \
		while [ $$CHECK_COUNT -lt 120 ]; do \
			CHECKS=$$(gh pr view "$$PR_NUMBER" --json statusCheckRollup -q ".statusCheckRollup | map(select(.status != \"COMPLETED\")) | length" 2>/dev/null); \
			if [ "$$CHECKS" = "0" ]; then \
				FAILED=$$(gh pr view "$$PR_NUMBER" --json statusCheckRollup -q ".statusCheckRollup | map(select(.conclusion == \"FAILURE\")) | length" 2>/dev/null); \
				if [ "$$FAILED" -gt 0 ]; then echo "ERROR: CI checks failed"; echo "Check results: $$PR_URL"; exit 1; fi; \
				echo "✓ All checks passed!"; \
				break; \
			fi; \
			sleep 1; CHECK_COUNT=$$((CHECK_COUNT + 1)); \
		done; \
		if [ $$CHECK_COUNT -ge 120 ]; then echo "ERROR: Checks did not complete within 2 minutes"; echo "Check status: $$PR_URL"; exit 1; fi; \
		echo ""; \
		echo "Step 7/10: Merging pull request..."; \
		gh pr merge --squash "$$PR_NUMBER" >/dev/null 2>&1 || { echo "ERROR: Failed to merge PR"; exit 1; }; \
		echo "✓ PR merged (squash)"; \
		echo ""; \
		echo "Step 8/10: Waiting for merge to complete..."; \
		MERGE_COUNT=0; \
		while [ $$MERGE_COUNT -lt 30 ]; do \
			STATE=$$(gh pr view "$$PR_NUMBER" --json state -q .state 2>/dev/null); \
			if [ "$$STATE" = "MERGED" ]; then echo "✓ Merge completed"; break; fi; \
			sleep 1; MERGE_COUNT=$$((MERGE_COUNT + 1)); \
		done; \
		echo ""; \
		echo "Step 9/10: Confirming tag creation..."; \
		TAG_COUNT=0; \
		while [ $$TAG_COUNT -lt 30 ]; do \
			if git rev-parse -q --verify "refs/tags/$$VERSION" >/dev/null 2>&1; then echo "✓ Tag $$VERSION created"; break; fi; \
			git fetch --tags origin >/dev/null 2>&1; sleep 1; TAG_COUNT=$$((TAG_COUNT + 1)); \
		done; \
		echo ""; \
		echo "Step 10/10: Cleaning up..."; \
		git checkout main >/dev/null 2>&1; \
		git branch -D "$$BRANCH" >/dev/null 2>&1; \
		git pull origin main >/dev/null 2>&1; \
		echo "✓ Deleted local branch $$BRANCH"; \
		echo ""; \
		echo "✓ Release $$VERSION created successfully!"; \
		echo ""; \
		echo "GitHub Actions will now build and create artifacts."; \
		echo "Monitor at: https://github.com/cutehax0r/proj/actions/workflows/release.yml"; \
		echo "Or run: gh run list --workflow release.yml -L 5"; \
	fi'
	@true

# Catch-all for release targets to suppress "No rule to make target" errors
%:
	@true
