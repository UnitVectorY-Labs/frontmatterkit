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
frontmatterkit validate post.md
```

### Validate from stdin

```bash
cat post.md | frontmatterkit validate
```

## Get

### Get all front matter as YAML

```bash
frontmatterkit get post.md
```

### Get a specific field

```bash
frontmatterkit get --path .title post.md
```

### Get a nested value

```bash
frontmatterkit get --path .author.name post.md
```

### Get an array element

```bash
frontmatterkit get --path '.tags[0]' post.md
```

### Get all front matter as JSON

```bash
frontmatterkit get --format json post.md
```

### Get a field as JSON

```bash
frontmatterkit get --format json --path .tags post.md
```

## Set

### Set a value

```bash
frontmatterkit set --set '.title=New Title' post.md
```

### Set multiple values

```bash
frontmatterkit set --set '.title=New Title' --set '.draft=true' post.md
```

### Set and write in place

```bash
frontmatterkit set --set '.title=New Title' --in-place post.md
```

### Set on a file without front matter

```bash
frontmatterkit set --set '.title=Hello' plain.md
```

## Unset

### Remove a field

```bash
frontmatterkit unset --path .draft post.md
```

### Remove a nested field

```bash
frontmatterkit unset --path .author.name post.md
```

## Assert

### Check that a field exists

```bash
frontmatterkit assert --assert '.title exists' post.md
```

### Check a value

```bash
frontmatterkit assert --assert '.draft == false' post.md
```

### Numeric comparison

```bash
frontmatterkit assert --assert '.count >= 1' post.md
```

### Array contains

```bash
frontmatterkit assert --assert '.tags contains "go"' post.md
```

### Multiple assertions

```bash
frontmatterkit assert --assert '.title exists' --assert '.draft == false' post.md
```

### Check that a field does not exist

```bash
frontmatterkit assert --assert '.obsolete not exists' post.md
```

## Pipelines

### Validate before committing

```bash
frontmatterkit validate post.md && git add post.md
```

### Extract a title for use in scripts

```bash
TITLE=$(frontmatterkit get --path .title post.md)
```

### Check conditions in CI

```bash
frontmatterkit assert --assert '.draft == false' --assert '.title exists' post.md
```