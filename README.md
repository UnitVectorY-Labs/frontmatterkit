# frontmatterkit

A Unix-style CLI for validating, querying, asserting, and minimally updating front matter in Markdown files.

[![Build](https://github.com/UnitVectorY-Labs/frontmatterkit/actions/workflows/build-go.yml/badge.svg)](https://github.com/UnitVectorY-Labs/frontmatterkit/actions/workflows/build-go.yml)
[![License: MIT](https://img.shields.io/github/license/UnitVectorY-Labs/frontmatterkit)](https://github.com/UnitVectorY-Labs/frontmatterkit/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/UnitVectorY-Labs/frontmatterkit)](https://go.dev/)

## Features

- **Validate** YAML front matter for correctness
- **Get** front matter values by path, output as YAML or JSON
- **Set** front matter fields with minimal document changes
- **Unset** front matter fields cleanly
- **Assert** conditions on front matter for CI/CD checks
- Reads from files or stdin — fits naturally into Unix pipelines
- Strict YAML validation with helpful exit codes

## Installation

```bash
go install github.com/UnitVectorY-Labs/frontmatterkit@latest
```

## Quick Usage

```bash
# Validate front matter
frontmatterkit validate post.md

# Get a field value
frontmatterkit get --path .title post.md

# Set a field
frontmatterkit set --set '.title=New Title' --in-place post.md

# Remove a field
frontmatterkit unset --path .draft --in-place post.md

# Assert conditions in CI
frontmatterkit assert --assert '.draft == false' --assert '.title exists' post.md
```

## Documentation

Full documentation is available at [unitvectory-labs.github.io/frontmatterkit](https://unitvectory-labs.github.io/frontmatterkit/).

## License

This project is licensed under the [MIT License](LICENSE).
