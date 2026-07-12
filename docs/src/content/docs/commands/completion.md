---
title: completion
description: Generate shell completion scripts
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
go-cli-docs completion [bash|zsh|fish|powershell]
```

## Source

See [root.go](https://github.com/imdevan/go-cli-docs/blob/main/cmd/go-cli-docs/root.go) for implementation details.
