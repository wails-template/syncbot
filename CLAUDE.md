# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go sync bot that clones Vite templates from `github.com/vitejs/vite`, copies them into the `wails-template` organization with Wails-specific template files, and commits/pushes changes. It runs on a schedule (every Monday at 10PM) via GitHub Actions.

## Commands

```bash
go run .  # Run the sync bot (requires GH_TOKEN env var)
go build  # Build the binary
```

## Architecture

**Main sync flow** (in `sync_vite.go`):
1. Clone the upstream Vite repo to a temp directory
2. Find all `template-*` directories in `packages/create-vite`
3. For each template: clone the target wails-template repo, clean all files, copy Wails template files, copy upstream frontend files, commit, and push

**Key modules:**
- `git.go` - `Repository` struct wrapping go-git operations (Clone, AddAll, Commit, Push, Status)
- `config.go` - Constants for GitHub username, email, target org URL, and `GH_TOKEN` env var
- `fs.go` - `Clean()` helper that removes all files/dirs except specified ignores
- `errors.go` - `ErrNoTemplatesFound` error

**Template directory** (`template/`):
- `app.tmpl.go`, `main.go.tmpl`, `go.tmpl.mod` - Go files copied into each synced repo
- `wails.tmpl.json` - Wails configuration template
- `template.json` - Template metadata
- `frontend/`, `scripts/` - Frontend and build scripts

**GitHub token**: Uses a GitHub App token via `actions/create-github-app-token` (not a PAT).

## CI

The sync runs on a cron schedule (`0 22 * * 1`) and also on push/PR to main. It requires `GH_TOKEN` to be set via environment variable.