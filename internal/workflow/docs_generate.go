// internal/workflow/docs_generate.go
package workflow

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// PackageInfo holds values from internal/package/package.toml
type PackageInfo struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    DocsSite    string `json:"docs_site"`
    DocsBase    string `json:"docs_base"`
    Repository  string `json:"repository"`
}

// CommandInfo mirrors the JSON produced by parse_commands.go
type CommandInfo struct {
    CmdName string `json:"cmd_name"`
    GoFile  string `json:"go_file"`
    Short   string `json:"short"`
    Use     string `json:"use"`
    Doc     string `json:"doc"`
    FlagGroups []struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        Flags []struct {
            Name        string `json:"name"`
            Short       string `json:"short"`
            Type        string `json:"type"`
            Description string `json:"description"`
        } `json:"flags"`
    } `json:"flag_groups"`
}

// GenerateDocs performs the full documentation generation using Go logic.
func GenerateDocs() error {
    // 1. Load package metadata
    pkgInfo, err := loadPackageInfo()
    if err != nil {
        return fmt.Errorf("failed to load package info: %w", err)
    }

    // 2. Detect commands via parse_commands.go (already used in InitAndGenerate)
    commandsJSON, err := exec.Command("go", "run", "./scripts/parse_commands.go", "./cmd/"+pkgInfo.Name).Output()
    if err != nil {
        return fmt.Errorf("failed to run parse_commands.go: %w", err)
    }
    var commands []CommandInfo
    if err := json.Unmarshal(commandsJSON, &commands); err != nil {
        return fmt.Errorf("failed to unmarshal commands json: %w", err)
    }

    // 3. Update docs/config.mjs
    if err := writeConfigMjs(pkgInfo); err != nil {
        return err
    }

    // 4. Generate sidebar.mjs
    if err := writeSidebarMjs(pkgInfo, commands); err != nil {
        return err
    }

    // 5. Generate content pages (index, install, configuration, contributing)
    if err := generateContentPage("README.md", "index.md", pkgInfo.Name, pkgInfo.Description); err != nil {
        return err
    }
    if err := generateContentPage("INSTALL.md", "install.md", "Install", fmt.Sprintf("Installation instructions for %s", pkgInfo.Name)); err != nil {
        return err
    }
    if err := generateContentPage("CONFIG.md", "configuration.md", "Configuration", fmt.Sprintf("Configuration options for %s", pkgInfo.Name)); err != nil {
        return err
    }
    // CONTRIBUTING.md is optional
    _ = generateContentPage("CONTRIBUTING.md", "contributing.md", "Contributing", fmt.Sprintf("Contributing to %s", pkgInfo.Name))

    // 6. Generate command pages
    if err := generateCommandPages(pkgInfo, commands); err != nil {
        return err
    }

    // 7. Generate API docs using gomarkdoc (if installed)
    if err := generateAPIDocs(pkgInfo); err != nil {
        return err
    }

    fmt.Println("✅ Documentation generation completed successfully!")
    return nil
}

func loadPackageInfo() (*PackageInfo, error) {
    pkgFile := filepath.Join("internal", "package", "package.toml")
    data, err := ioutil.ReadFile(pkgFile)
    if err != nil {
        return nil, err
    }
    info := &PackageInfo{}
    scanner := bufio.NewScanner(bytes.NewReader(data))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        key := strings.TrimSpace(parts[0])
        val := strings.Trim(strings.TrimSpace(parts[1]), "\"")
        switch key {
        case "name":
            info.Name = val
        case "description":
            info.Description = val
        case "docs_site":
            info.DocsSite = val
        case "docs_base":
            info.DocsBase = val
        case "repository":
            info.Repository = val
        }
    }
    if info.Name == "" {
        // fallback to project directory name
        cwd, _ := os.Getwd()
        info.Name = filepath.Base(cwd)
    }
    return info, nil
}

func writeConfigMjs(info *PackageInfo) error {
    path := filepath.Join("docs", "config.mjs")
    content := fmt.Sprintf(`const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"

export default {
  url: isProduction ? "%s" : "http://localhost:4321",
  basePath: isProduction ? "%s" : "/",
  github: "%s",
  githubDocs: "%s",
  title: "%s",
  description: "%s",
}
`, info.DocsSite, info.DocsBase, info.Repository, info.Repository, info.Name, info.Description)
    return os.WriteFile(path, []byte(content), 0644)
}

func writeSidebarMjs(info *PackageInfo, commands []CommandInfo) error {
    // Build commands list for sidebar
    var cmdItems []string
    for _, c := range commands {
        if c.CmdName == info.Name {
            continue // skip root
        }
        cmdItems = append(cmdItems, fmt.Sprintf("    { label: '%s', link: '/commands/%s' }", c.CmdName, c.CmdName))
    }
    commandsSection := strings.Join(cmdItems, ",\n")

    // API packages detection (internal folders with .go files)
    apiPackages, err := detectAPIPackages()
    if err != nil {
        return err
    }
    var apiItems []string
    for _, p := range apiPackages {
        apiItems = append(apiItems, fmt.Sprintf("    { label: '%s', link: '/api/%s' }", p, p))
    }
    apiSection := strings.Join(apiItems, ",\n")

    content := fmt.Sprintf(`import config from './config.mjs';

const apiReference = {
  label: 'API Reference',
  items: [
%s    {
      label: 'Adapters',
      items: [],
    },
  ],
};

const sidebar = [
  { label: '%s', link: '/' },
  { label: 'Install', link: '/install' },
  { label: 'Commands', items: [
    { label: '%s', link: '/commands/%s' },
%s  ] },
  { label: 'Configuration', link: '/configuration' },
];

const isProduction = process.env.NODE_ENV === 'production';
if (!isProduction) {
  sidebar.push(apiReference);
}

export default sidebar;
`, apiSection, info.Name, info.Name, info.Name, commandsSection)
    path := filepath.Join("docs", "sidebar.mjs")
    return os.WriteFile(path, []byte(content), 0644)
}

func detectAPIPackages() ([]string, error) {
    internalPath := "internal"
    entries, err := os.ReadDir(internalPath)
    if err != nil {
        return nil, err
    }
    var pkgs []string
    for _, e := range entries {
        if !e.IsDir() {
            continue
        }
        name := e.Name()
        // skip special dirs
        if name == "testutil" || name == "adapters" || name == "package" {
            continue
        }
        // check for .go files
        files, _ := os.ReadDir(filepath.Join(internalPath, name))
        for _, f := range files {
            if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
                pkgs = append(pkgs, name)
                break
            }
        }
    }
    return pkgs, nil
}

func generateContentPage(srcFile, dstFile, title, description string) error {
    if _, err := os.Stat(srcFile); err != nil {
        return nil // optional source missing
    }
    data, err := ioutil.ReadFile(srcFile)
    if err != nil {
        return err
    }
    outDir := filepath.Join("docs", "src", "content", "docs")
    if err := os.MkdirAll(outDir, 0755); err != nil {
        return err
    }
    dstPath := filepath.Join(outDir, dstFile)
    var buf bytes.Buffer
    buf.WriteString("---\n")
    buf.WriteString(fmt.Sprintf("title: %s\n", title))
    buf.WriteString(fmt.Sprintf("description: %s\n", description))
    buf.WriteString("---\n\n")
    buf.Write(data)
    return os.WriteFile(dstPath, buf.Bytes(), 0644)
}

func generateCommandPages(pkgInfo *PackageInfo, commands []CommandInfo) error {
    outDir := filepath.Join("docs", "src", "content", "docs", "commands")
    if err := os.MkdirAll(outDir, 0755); err != nil {
        return err
    }
    // root command page
    var rootCmd *CommandInfo
    for _, c := range commands {
        if c.CmdName == pkgInfo.Name {
            rootCmd = &c
            break
        }
    }
    if rootCmd != nil {
        if err := writeCommandPage(outDir, pkgInfo.Name, rootCmd, pkgInfo.Description); err != nil {
            return err
        }
    }
    // subcommands
    for _, c := range commands {
        if c.CmdName == pkgInfo.Name {
            continue
        }
        if err := writeCommandPage(outDir, c.CmdName, &c, ""); err != nil {
            return err
        }
    }
    return nil
}

func writeCommandPage(dir, cmdName string, info *CommandInfo, fallbackDesc string) error {
    dstPath := filepath.Join(dir, fmt.Sprintf("%s.md", cmdName))
    var buf bytes.Buffer
    title := cmdName
    desc := info.Doc
    if desc == "" || desc == "null" {
        desc = fallbackDesc
    }
    buf.WriteString("---\n")
    buf.WriteString(fmt.Sprintf("title: %s\n", title))
    buf.WriteString(fmt.Sprintf("description: %s\n", desc))
    buf.WriteString("---\n\n")
    buf.WriteString(desc + "\n\n")
    // Usage
	buf.WriteString("## Usage\n\n```go\n")
	usage := info.Use
	if usage == "" || usage == "null" {
		usage = fmt.Sprintf("%s %s", cmdName, "")
	}
	buf.WriteString(usage + "\n```\n\n")
    // Flags
    if len(info.FlagGroups) > 0 {
        buf.WriteString("## Flags\n\n")
        for _, fg := range info.FlagGroups {
            buf.WriteString(fmt.Sprintf("### %s\n\n", fg.Name))
            if fg.Description != "" && fg.Description != "null" {
                buf.WriteString(fg.Description + "\n\n")
            }
            if len(fg.Flags) > 0 {
                buf.WriteString("| Flag | Type | Description |\n|------|------|-------------|\n")
                for _, fl := range fg.Flags {
                    flagCol := "--" + fl.Name
                    if fl.Short != "" && fl.Short != "null" {
                        flagCol = fmt.Sprintf("-%s, --%s", fl.Short, fl.Name)
                    }
                    typ := strings.ToLower(fl.Type)
                    descr := fl.Description
                    buf.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", flagCol, typ, descr))
                }
                buf.WriteString("\n")
            }
        }
    }
    return os.WriteFile(dstPath, buf.Bytes(), 0644)
}

func generateAPIDocs(pkgInfo *PackageInfo) error {
    // Ensure gomarkdoc is available
    if _, err := exec.LookPath("gomarkdoc"); err != nil {
        fmt.Println("Installing gomarkdoc...")
        install := exec.Command("go", "install", "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest")
        install.Stdout = os.Stdout
        install.Stderr = os.Stderr
        if err := install.Run(); err != nil {
            return fmt.Errorf("failed to install gomarkdoc: %w", err)
        }
    }
    apiDir := filepath.Join("docs", "src", "content", "docs", "api")
    tmpDir, err := ioutil.TempDir("", "api_tmp")
    if err != nil {
        return err
    }
    defer os.RemoveAll(tmpDir)

    pkgs, err := detectAPIPackages()
    if err != nil {
        return err
    }
    for _, pkg := range pkgs {
        pkgPath := filepath.Join("internal", pkg)
        outFile := filepath.Join(tmpDir, pkg+".md")
        cmd := exec.Command("gomarkdoc", "--output", outFile, "--template-file", "file=./docs/templates/file.gotxt", pkgPath)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            fmt.Printf("⚠️  No exported symbols in %s\n", pkg)
            continue
        }
        content, _ := os.ReadFile(outFile)
        final := fmt.Sprintf("---\ntitle: %s\ndescription: API documentation for the %s package\n---\n\n%s", pkg, pkg, string(content))
        dstPath := filepath.Join(apiDir, pkg+".md")
        os.MkdirAll(apiDir, 0755)
        os.WriteFile(dstPath, []byte(final), 0644)
    }
    // adapters
    adaptersPath := filepath.Join("internal", "adapters")
    if entries, _ := os.ReadDir(adaptersPath); len(entries) > 0 {
        adaptersOut := filepath.Join(apiDir, "adapters")
        os.MkdirAll(adaptersOut, 0755)
        for _, a := range entries {
            if !a.IsDir() {
                continue
            }
            name := a.Name()
            pkgPath := filepath.Join(adaptersPath, name)
            outFile := filepath.Join(tmpDir, name+".md")
            cmd := exec.Command("gomarkdoc", "--output", outFile, "--template-file", "file=./docs/templates/file.gotxt", pkgPath)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            if err := cmd.Run(); err != nil {
                fmt.Printf("⚠️  No exported symbols in adapter %s\n", name)
                continue
            }
            content, _ := os.ReadFile(outFile)
            final := fmt.Sprintf("---\ntitle: adapters/%s\ndescription: API documentation for the %s adapter\n---\n\n%s", name, name, string(content))
            dstPath := filepath.Join(adaptersOut, name+".md")
            os.WriteFile(dstPath, []byte(final), 0644)
        }
    }
    return nil
}
