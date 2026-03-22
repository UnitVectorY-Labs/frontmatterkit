---
layout: default
title: Examples
nav_order: 4
permalink: /examples
---

# Examples
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

- TOC
{:toc}

---

These examples are intentionally few and complete. Each one is backed by a data-driven CLI case under `testdata/`, and each one shows the Markdown input, the generated output, or both.
{: .note }

For the full flag reference, see the [Usage](/usage) page.
{: .note }

## Validate

### Reject malformed front matter

Backed by `testdata/validate/invalid-yaml`. Use this example when you want to see what a real validation failure looks like.
{: .warning }

```bash
frontmatterkit validate --in invalid-yaml.md
```

**Input Markdown**

```markdown
---
title: [broken
---
Body content.
```

**Error Output**

```text
Error: validation failed: invalid YAML in front matter: yaml: line 1: did not find expected ',' or ']'
```

This command exits with code `1`.

## Get

### Extract an array as JSON for another tool

Backed by `testdata/get/json-tags`. This is a good example of using `get` as a structured data extractor instead of just printing YAML back out.
{: .important }

```bash
frontmatterkit get --format json --path .tags --in json-tags.md
```

**Input Markdown**

```markdown
---
title: Hello World
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

**Output**

```json
[
  "go",
  "yaml"
]
```

## Set

### Update front matter from stdin and emit the rewritten document

Backed by `testdata/set/stdin-set`. This shows the Unix-friendly stdin workflow while still giving a complete before-and-after example.
{: .note }

```bash
cat stdin-set.md | frontmatterkit set --set '.title=Stdin Title'
```

**Input Markdown**

```markdown
---
title: Hello World
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

**Output Markdown**

```markdown
---
title: Stdin Title
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

## Unset

### Remove a publishing flag and write a clean output file

Backed by `testdata/unset/out-file`. This is the useful `unset` workflow: remove one field, preserve everything else, and write the result to a separate file.
{: .important }

```bash
frontmatterkit unset --path .draft --in out-file.md --out out-file.out.md
```

**Input Markdown**

```markdown
---
title: Hello World
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

**Output Markdown**

```markdown
---
title: Hello World
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

## Assert

### Fail fast when required front matter is missing

Backed by `testdata/assert/exists-fail`. This is the most useful `assert` pattern for CI: enforce a required field and stop the pipeline with a clear failure.
{: .warning }

```bash
frontmatterkit assert --assert '.missing exists' --in exists-fail.md
```

**Input Markdown**

```markdown
---
title: Hello World
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test User
---
This is the body content.
```

**Error Output**

```text
Error: assertion ".missing exists" failed: assertion failed: .missing does not exist
```

This command exits with code `1`.
