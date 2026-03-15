# HelmWatch | Watch - Diff - Deploy

**HelmWatch** is an interactive CLI tool for watching and diffing Helm charts in real-time. It allows you to see changes live as you modify your `values.yaml` or chart templates with a terminal UI, search and colorized diffs.

![HelmWatch Screenshot](./assets/helmwatch.svg)

## Features

- 🌟 **On-Change Diff**: Automatically detects changes in Helm templates and values files and shows the diff of the generated manifests.
- 🖥 **Interactive TUI**: Browse diffs, search and rerender charts directly from the terminal.  
- 🎨 **Colorized diffs**

---

## Installation

```bash
go install github.com/eklatzer/helmwatch@latest
```

## Requirements

- [Helm](https://helm.sh/docs/intro/install/)
- `diff`

## Usage

```bash
helmwatch --help   
Interactive Helm diff watcher

Usage:
  helmwatch [flags]
  helmwatch [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version

Flags:
      --chart string     Path to the chart or remote chart reference (default ".")
      --config string    Path to the helmwatch config (default "helmwatch.yaml")
  -h, --help             help for helmwatch
      --values string    Path to the values file (default "values.yaml")
      --version string   Version of the chart
```

### Example Usage

```bash
# watch changes of local helm chart
helmwatch

# watch changes using a remote helm chart with a local helm values file
# requires helm repo add ...
helmwatch --chart argo/argo-cd --version 9.3.7
```

### Configuration

```yaml
# exclusions is a list of substrings
# any line containing a exclusion, is excluded in the diff
exclusions: []
```

Example:
```yaml
exclusions:
  - "checksum/secret"
  - "checksum/config"
```