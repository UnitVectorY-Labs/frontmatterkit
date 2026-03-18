---
layout: default
title: Usage
nav_order: 3
permalink: /usage
---

# Usage
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

- TOC
{:toc}

---

## Global Conventions

### Input

Every subcommand accepts an optional file argument. If no file is provided, input is read from stdin. This makes frontmatterkit composable with other Unix tools via pipes.

```
frontmatterkit <subcommand> [flags] [file]
```

### Exit Codes

frontmatterkit uses structured exit codes so scripts and CI pipelines can react appropriately.

| Code | Meaning |
|------|---------|
| `0`  | Success |
| `1`  | Validation or assertion failure |
| `2`  | Usage or syntax error |
| `3`  | I/O error |

---

## validate

Check that the YAML front matter in a Markdown file is well-formed.

### Syntax

```
frontmatterkit validate [file]
```

### Flags

This subcommand has no additional flags.

### Behavior

- Reads from stdin if no file argument is provided.
- A file with no front matter is considered **valid**.
- If the file begins with `---` and the YAML content between the delimiters is malformed, validation **fails** with exit code `1`.

---

## get

Extract front matter values from a Markdown file.

### Syntax

```
frontmatterkit get [--format yaml|json] [--path <path>] [file]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `yaml` | Output format: `yaml` or `json`. |
| `--path` | `.` | jq-like path to the value to extract. Use `.` for the full front matter object. |

### Behavior

- Missing front matter is treated as an empty object (`{}`).
- If the specified path does not exist, the output is `null`.

---

## set

Add or update front matter fields in a Markdown file.

### Syntax

```
frontmatterkit set [--set .path=value]... [--from <file>] [--mode overwrite|patch] [--in-place] [--output <path>] [file]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--set` | | Set a value using `.path=yamlValue` format. Can be specified multiple times. |
| `--from` | | Read YAML values from a file. Use `-` for stdin. |
| `--mode` | `overwrite` | Merge strategy: `overwrite` replaces the entire front matter, `patch` merges values into the existing front matter. |
| `--in-place` | `false` | Overwrite the source file with the result. |
| `--output` | | Write the result to the specified file path instead of stdout. |

### Behavior

- If the file has no front matter, a new front matter block is created.
- Values provided via `--set` are interpreted as YAML, so `true`, `1`, and `[a, b]` produce their respective YAML types rather than plain strings.
- The document body outside of the front matter is preserved unchanged.

---

## unset

Remove a field from the front matter of a Markdown file.

### Syntax

```
frontmatterkit unset --path <path> [--in-place] [--output <path>] [file]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | | Path to the field to remove. **Required.** |
| `--in-place` | `false` | Overwrite the source file with the result. |
| `--output` | | Write the result to the specified file path instead of stdout. |

### Behavior

- If the specified path does not exist in the front matter, the command is a **no-op** — the file is output unchanged.
- The document body outside of the front matter is preserved unchanged.

---

## assert

Test conditions against front matter fields. Useful for CI checks and pre-commit hooks.

### Syntax

```
frontmatterkit assert --assert '<expr>'... [file]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--assert` | | Assertion expression to evaluate. Can be specified multiple times. **Required.** |

### Operators

| Operator | Description |
|----------|-------------|
| `exists` | Field exists (e.g., `.title exists`) |
| `not exists` | Field does not exist (e.g., `.draft not exists`) |
| `==` | Equal to |
| `!=` | Not equal to |
| `<` | Less than |
| `<=` | Less than or equal to |
| `>` | Greater than |
| `>=` | Greater than or equal to |
| `contains` | Array contains a value (e.g., `.tags contains "go"`) |
| `not contains` | Array does not contain a value |

### Behavior

- When multiple `--assert` flags are provided, all assertions must pass (logical AND). If any assertion fails, the command exits with code `1`.

---

## Path Syntax

Paths use a jq-like dot notation to address values within the front matter YAML structure. Paths are used by the `get`, `set`, `unset`, and `assert` subcommands.

| Path | Description |
|------|-------------|
| `.` | Root — the full front matter object |
| `.field` | Top-level field |
| `.a.b.c` | Nested fields |
| `.tags[0]` | Array element by index |
| `.a.b[2].c` | Mixed nesting with array indices |