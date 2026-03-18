---
layout: default
title: frontmatterkit
nav_order: 1
permalink: /
---

# frontmatterkit

A Unix-style CLI for validating, querying, asserting, and minimally updating front matter in Markdown files.

---

## What is frontmatterkit?

**frontmatterkit** is a command-line tool purpose-built for working with YAML front matter in Markdown files. Whether you're enforcing content standards in CI, scripting bulk edits across a documentation site, or just querying metadata from a blog post, frontmatterkit gives you precise, composable commands that fit naturally into Unix pipelines.

## Key Features

- **Validate** — Verify that YAML front matter is well-formed, catching syntax errors before they break your build.
- **Get** — Extract front matter values by path, with output in YAML or JSON for easy downstream processing.
- **Set** — Add or update front matter fields while preserving the rest of your document exactly as-is.
- **Unset** — Remove fields cleanly without touching anything else.
- **Assert** — Test conditions on front matter fields with expressive operators — ideal for CI/CD gates and pre-commit hooks.

## Why frontmatterkit?

- **Strict validation** — Malformed YAML is caught immediately with clear exit codes, not silently ignored.
- **Minimal changes** — The `set` and `unset` commands modify only what you ask for; the rest of the file is left byte-for-byte identical.
- **Unix-friendly** — Commands read from `--in` or stdin, write to `--out` or stdout, and use meaningful exit codes. Compose them with `grep`, `jq`, `xargs`, or anything else in your toolkit.
- **CI-ready assertions** — Replace fragile `grep` checks with structured assertions like `.draft == false` and `.tags contains "go"`.

## Quick Start

```bash
# Install
go install github.com/UnitVectorY-Labs/frontmatterkit@latest

# Validate a Markdown file
frontmatterkit validate --in post.md

# Get the title
frontmatterkit get --path .title --in post.md

# Ensure the post is not a draft and has a title
frontmatterkit assert --assert '.draft == false' --assert '.title exists' --in post.md

# Show command-specific help
frontmatterkit set help
```

## Documentation

- [Usage Reference]({% link USAGE.md %}) — Complete CLI reference with all subcommands and flags.
- [Examples]({% link EXAMPLES.md %}) — Practical examples for common workflows.
