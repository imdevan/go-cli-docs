---
title: init
description: Scaffold the Astro Starlight docs directory
---

Creates the docs/ folder and scaffolds Astro Starlight.
If the docs/ directory already exists, init is a no-op.
After scaffolding, init automatically runs generate to populate the docs site.

### Example

```bash
go-cli-docs init
go-cli-docs init --pkg-manager npm
```

## Usage

```bash
go-cli-docs init
```

## Flags

### Flags



| Flag | Type | Description |
|------|------|-------------|
| -p, --pkg-manager | string | Package manager to use (bun, npm, yarn, pnpm) |

### Global Flags



| Flag | Type | Description |
|------|------|-------------|
| -a, --gen-api-docs | bool | Generate API documentation via gomarkdoc |
| -t, --templates | stringarray | Path to a file or directory of custom templates overriding the embedded defaults (repeatable) |


## Source

See [init.go](https://github.com/imdevan/go-cli-docs/blob/main/cmd/go-cli-docs/init.go) for implementation details.
