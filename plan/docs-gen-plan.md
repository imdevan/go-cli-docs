# Docs Generation Refactor Plan (Go)

## Overview

This document outlines a full refactor of a bash-based documentation generation script into a structured, modular Go system using templates.

The goal is to replace:
- jq parsing
- bash loops
- heredocs
- sed hacks

With:
- typed structs
- Go templates
- modular packages
- testable components

---

current bash script: /home/devy/Projects/clis/go-cli-docs/scripts/docs_generate.sh


## goal

This project aims to build this project /home/devy/Projects/clis/go-cli-docs

Into a go cli package that can be used to generate docs for go cli packages. 

It is self documenting in that generating the docs for this project will mean that it can 
successfully build docs for other projects. 

---

## Boilerplate setup

1. Remove internal/ui elements and deps. 
  - this project will be strictly go-cobra -> generate astro starlight docs.
2. Map commands to cmd/go-cli-docs/...
3. Map logic to internal/...


## Config variables

```toml
DOCS_DIR = "docs" # location of docs

```

## Flags

```bash
go-cli-docs -i, --init  init docs, or re-init
```


## Root command

go-cli-docs

Do the following: 
  if no DOCS_DIR:
    call init function
  else:
    cal main  generate function

## Primary functions

### init: install astro-starlight docs (version locked)
  create docs in ./docs folder
  then call generate function

### generate: generate docs

1. Read package metadata (TOML)
2. Generate markdown content pages
3. Parse Cobra commands
4. Generate command documentation
5. Generate API documentation via gomarkdoc 
  - optional based on gen_api_docs argument passed
    - in production defaults to false
    - in development defaults to true
6. Generate config
7. Generate sidebar files
    


## Templates

Templates should live in:

internal/templates/

Examples:
- config.mjs.tmpl
- sidebar.mjs.tmpl
- command.md.tmpl
- page.md.tmpl

---

## Embedding Templates

Use Go embed:

//go:embed *.tmpl
var FS embed.FS

---

## Template Example (config)

export default {
  url: isProduction ? "{{ .DocsSite }}" : "http://localhost:4321",
  basePath: isProduction ? "{{ .DocsBase }}" : "/",
  github: "{{ .Repository }}",
  title: "{{ .ProjectName }}",
  description: "{{ .Description }}",
}

---

## Template Example (sidebar)

{{ range .Commands }}
{ label: '{{ .Name }}', link: '/commands/{{ .Name }}' },
{{ end }}

---

## Template Example (command markdown)

---
title: {{ .Name }}
description: {{ .Description }}
---

{{ .Doc }}

## Usage

{{ .Usage }}

---

## Flag Groups Template

{{ if .FlagGroups }}
## Flags

{{ range .FlagGroups }}
### {{ .Name }}

{{ .Description }}

{{ if .Example }}
#### Example
{{ .Example }}
{{ end }}

| Flag | Type | Description |
|------|------|-------------|
{{ range .Flags }}
| {{ .Display }} | {{ .Type }} | {{ .Description }} |
{{ end }}

{{ end }}
{{ end }}

---

## Data Structures

type CommandDoc struct {
    Name        string
    Description string
    Usage       string
    Doc         string
    FlagGroups  []FlagGroup
}

type FlagGroup struct {
    Name        string
    Description string
    Example     string
    Flags       []Flag
}

type Flag struct {
    Display     string
    Type        string
    Description string
}

---

## Rendering Helper

func RenderToFile(tmpl *template.Template, path string, data any) error {
    f, _ := os.Create(path)
    defer f.Close()
    return tmpl.Execute(f, data)
}

---

## External Commands

Use exec.Command for gomarkdoc instead of bash.

---

## Design Principles

- Separate logic from presentation
- Prefer structs over string concatenation
- Use templates for all output
- Keep generators modular and composable

---

## Optional Enhancements

- Allow custom template directory override
- Add themes
- Add HTML output
- Add testing per generator

---

## Outcome

A maintainable, testable, extensible documentation generation system in Go.
