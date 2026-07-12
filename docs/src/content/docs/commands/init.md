---
title: init
description: Scaffold the Astro Starlight docs directory
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
go-cli-docs init
```

## Flags




| Flag | Type | Description |
|------|------|-------------|
| -p, --pkg-manager | string | Package manager to use (bun, npm, yarn, pnpm) |


## Source

See [root.go](https://github.com/imdevan/go-cli-docs/blob/main/cmd/go-cli-docs/root.go) for implementation details.
