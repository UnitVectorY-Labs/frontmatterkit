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

{: .note }
A test case must exist for each example in this file.

## Validate

### Validate a file

```bash
frontmatterkit validate --in post.md
```

### Validate from stdin

```bash
cat post.md | frontmatterkit validate
```

### Show validate help

```bash
frontmatterkit validate help
```

## Get

### Get all front matter as YAML

```bash
frontmatterkit get --in post.md
```

### Get a specific field

```bash
frontmatterkit get --path .title --in post.md
```

### Get a nested value

```bash
frontmatterkit get --path .author.name --in post.md
```

### Get an array element

```bash
frontmatterkit get --path '.tags[0]' --in post.md
```

### Get all front matter as JSON

```bash
frontmatterkit get --format json --in post.md
```

### Get a field as JSON

```bash
frontmatterkit get --format json --path .tags --in post.md
```

### Write extracted output to a file

```bash
frontmatterkit get --path .title --in post.md --out title.txt
```

### Show get help

```bash
frontmatterkit get help
```

## Set

### Set a value

```bash
frontmatterkit set --set '.title=New Title' --in post.md
```

### Set multiple values

```bash
frontmatterkit set --set '.title=New Title' --set '.draft=true' --in post.md
```

### Set and write in place

```bash
frontmatterkit set --set '.title=New Title' --in post.md --in-place
```

### Set on a file without front matter

```bash
frontmatterkit set --set '.title=Hello' --in plain.md
```

### Set and write to a separate file

```bash
frontmatterkit set --set '.title=Hello' --in post.md --out updated.md
```

### Show set help

```bash
frontmatterkit set help
```

## Unset

### Remove a field

```bash
frontmatterkit unset --path .draft --in post.md
```

### Remove a nested field

```bash
frontmatterkit unset --path .author.name --in post.md
```

### Remove a field in place

```bash
frontmatterkit unset --path .draft --in post.md --in-place
```

### Show unset help

```bash
frontmatterkit unset help
```

## Assert

### Check that a field exists

```bash
frontmatterkit assert --assert '.title exists' --in post.md
```

### Check a value

```bash
frontmatterkit assert --assert '.draft == false' --in post.md
```

### Numeric comparison

```bash
frontmatterkit assert --assert '.count >= 1' --in post.md
```

### Array contains

```bash
frontmatterkit assert --assert '.tags contains "go"' --in post.md
```

### Multiple assertions

```bash
frontmatterkit assert --assert '.title exists' --assert '.draft == false' --in post.md
```

### Check that a field does not exist

```bash
frontmatterkit assert --assert '.obsolete not exists' --in post.md
```

### Show assert help

```bash
frontmatterkit assert help
```

## Pipelines

### Validate before committing

```bash
frontmatterkit validate --in post.md && git add post.md
```

### Extract a title for use in scripts

```bash
TITLE=$(frontmatterkit get --path .title --in post.md)
```

### Check conditions in CI

```bash
frontmatterkit assert --assert '.draft == false' --assert '.title exists' --in post.md
```
