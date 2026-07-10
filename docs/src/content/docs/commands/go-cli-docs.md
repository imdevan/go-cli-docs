---
title: go-cli-docs
description: Generate Astro Starlight documentation for Go CLI projects
---

Generate Astro Starlight documentation for Go CLI projects.
The tool parses Cobra commands and flags, rendering markdown pages,
sidebar configs, and API docs.

### Example

```bash
go-cli-docs init
go-cli-docs generate
go-cli-docs watch
```

## Usage

```bash
go-cli-docs
```

## Flags

### Flags



| Flag | Type | Description |
|------|------|-------------|
| -v, --version | bool | Print version and exit |
| --config | string | Path to config file |


## Available Commands


- [`completion`](/commands/completion) - Generate shell completion scripts
- [`generate`](/commands/generate) - Generate all documentation from source
- [`init`](/commands/init) - Scaffold the Astro Starlight docs directory
- [`watch`](/commands/watch) - Watch source files and re-generate documentation on change

## Source

See [root.go](https://github.com/imdevan/go-cli-docs/blob/main/cmd/go-cli-docs/root.go) for implementation details.
