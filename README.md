[![GitHub release](https://img.shields.io/github/release/UnitVectorY-Labs/frontmatterkit.svg)](https://github.com/UnitVectorY-Labs/frontmatterkit/releases/latest) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Active](https://img.shields.io/badge/Status-Active-green)](https://guide.unitvectorylabs.com/bestpractices/status/#active) [![Go Report Card](https://goreportcard.com/badge/github.com/UnitVectorY-Labs/frontmatterkit)](https://goreportcard.com/report/github.com/UnitVectorY-Labs/frontmatterkit) 

# frontmatterkit

A Unix-style CLI for validating, querying, asserting, and minimally updating front matter in Markdown files.

## Features

- **Validate** YAML front matter for correctness
- **Get** front matter values by path, output as YAML or JSON
- **Set** front matter fields with minimal document changes
- **Unset** front matter fields cleanly
- **Assert** conditions on front matter for CI/CD checks
- Reads from `--in` or stdin and writes to `--out` or stdout
- Strict YAML validation with helpful exit codes

## Installation

```bash
go install github.com/UnitVectorY-Labs/frontmatterkit@latest
```

## Quick Usage

```bash
# Validate front matter
frontmatterkit validate --in post.md

# Get a field value
frontmatterkit get --path .title --in post.md

# Set a field
frontmatterkit set --set '.title=New Title' --in post.md --in-place

# Remove a field
frontmatterkit unset --path .draft --in post.md --in-place

# Assert conditions in CI
frontmatterkit assert --assert '.draft == false' --assert '.title exists' --in post.md

# Show command-specific help
frontmatterkit get help
```

## License

This project is licensed under the [MIT License](LICENSE).
