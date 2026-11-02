---
name: git-commit
description: Creates a git commit when the user asks to commit changes, to create a git
    commit, or mentions "/commit".
allowed-tools: Bash
---

# Git Commit with Conventional Commits

## Overview

Create standardized, semantic git commits using the Conventional Commits specification.
Analyze the actual diff to determine the appropriate type, scope, and message.

## Create Git Commit According to Project Structure

Structure the commit body by the project's directory structure (see the project README):
for each affected directory, list its changes under one of these sections:
- added
- modified
- removed

So, if an object (file or package) has several changes, first specify the package (its directory),
then the pattern of changes.

## Workflow

### 1. Analyze Diff

```bash
# If files are staged, use staged diff
git diff --staged

# If nothing staged, use working tree diff
git diff

# Also check status
git status --porcelain
```

### 2. Stage Files (if needed)

If nothing is staged or you want to group changes differently:

```bash
# Stage specific files
git add path/to/file1 path/to/file2

# Interactive staging
git add -p
```

**Never commit secrets** (.env, credentials.json, private keys).

### 3. Generate Commit Message

Analyze the diff to determine:

- **Type**: What kind of change is this?
- **Scope**: What area/module is affected?
- **Description**: One-line summary of what changed (present tense, imperative mood, <72 chars)

### 4. Execute Commit

```bash
# Single line
git commit -m "<type>[scope]: <description>"

# Multi-line with body/footer
git commit -m "$(cat <<'EOF'
<type>[scope]: <description>

<optional body>

<optional footer>
EOF
)"
```

### 5. Print the git commit message to the user
Template:

```text
======= git-commit-message =======
<commit message>
==================================
```

### 6. Push commit if user specifies it explicitly
```bash
# Push committed changes to the remote
git push
```

## Best Practices

- Imperative present tense: "add feature", not "added feature"

## Git Safety Protocol

- NEVER update git config
- NEVER run destructive commands (--force, hard reset) without an explicit request
- NEVER skip hooks (--no-verify) unless the user asks
- NEVER force push to main/master
- If the commit fails due to hooks, fix and create a NEW commit (don't amend)