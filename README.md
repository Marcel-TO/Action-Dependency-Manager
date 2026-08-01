# Action Dependency Manager (ADM)

A small command-line tool that helps you keep the GitHub Action dependencies in
your workflow files up to date — across many local repositories at once.

ADM scans the `.github/workflows` directory of each configured repository,
discovers every `uses:` reference, resolves the latest release of each action
and can either **report** what is outdated (`check`) or **pin** every action to
the latest release commit (`update`).

## Why?

Pinning actions to a commit SHA is the recommended way to consume third-party
GitHub Actions securely, but a plain SHA tells you nothing about which version
it represents, and keeping it current by hand across several repositories is
tedious. ADM automates this by rewriting references to the following form:

```yaml
uses: actions/checkout@3d3c42e... # v7.0.1
```

The commit SHA is used as the actual pin, while the human-readable release tag
is kept as a trailing comment.

## Features

- 🔍 Discovers all workflow files under `.github/workflows` in each repository.
- 🧩 Parses every `uses:` reference and extracts the dependency and its version.
- 🌐 Resolves the latest release tag and its commit SHA via the `gh` CLI.
- ✅ `check` — read-only report of which dependencies are outdated, grouped per
  workflow file.
- ⬆️ `update` — pins every dependency to its latest release commit SHA with the
  tag preserved as a comment.
- 📁 Works across multiple repositories from a single config file, with an
  optional `--select` filter.

## Prerequisites

ADM relies on the [GitHub CLI](https://cli.github.com/) to look up release tags
and their corresponding commit SHAs.

- [Go](https://go.dev/dl/) 1.26+ (to build/run from source)
- [`gh` CLI](https://cli.github.com/) — installed and authenticated
  (`gh auth login`)
- A GitHub account with access to the actions you depend on

Behind the scenes ADM runs commands equivalent to:

```bash
# Fetch the latest release tag of a dependency
gh release view --repo actions/checkout --json tagName --jq .tagName

# Resolve a tag (lightweight or annotated) to its commit SHA
gh api repos/actions/checkout/commits/v7.0.1 -H "Accept: application/vnd.github.sha"
```

## Configuration

ADM reads a `action-repositories.yaml` file from the current working directory.
Copy the provided example to get started:

```bash
cp action-repositories.yaml.example action-repositories.yaml
```

```yaml
repos:
  - name: "Repo 1"
    id: 1
    repoLocation: "local/path/to/repo1"
  - name: "Repo 2"
    id: 2
    repoLocation: "local/path/to/repo2"
```

| Field          | Description                                                        |
| -------------- | ------------------------------------------------------------------ |
| `name`         | Human-friendly name shown in the output.                           |
| `id`           | Numeric identifier used by the `--select` flag.                    |
| `repoLocation` | Absolute or relative path to the local checkout of the repository. |

> `action-repositories.yaml` is git-ignored, so your local paths stay private.

## Usage

Run directly with Go:

```bash
go run main.go <command> [flags]
```

Or build a binary first:

```bash
go build -o adm .
./adm <command> [flags]
```

### Commands

| Command  | Description                                                              |
| -------- | ----------------------------------------------------------------------- |
| `check`  | Reports which dependencies have a newer version. Does **not** modify files. |
| `update` | Pins every dependency to its latest release commit SHA in the workflow files. |

### Flags

| Flag       | Commands        | Default | Description                                                                    |
| ---------- | --------------- | ------- | ------------------------------------------------------------------------------ |
| `--select` | `check`, `update` | *(all)* | Comma-separated list of repository `id`s to process. Empty means all of them. |
| `--debug`  | `check`, `update` | `false` | Enable verbose debug logging.                                                  |
| `--commit` | `update`        | `true`  | Whether to commit the changes after updating.                                  |

### Examples

```bash
# Report outdated dependencies across every configured repository
go run main.go check

# Only check the repositories with id 1 and 3
go run main.go check --select="1,3"

# Update dependencies in every repository
go run main.go update

# Update only repository 1, with verbose logging
go run main.go update --select="1" --debug
```

### Example output

```
azure-fa-jaw-mobilfunk.yml
--------------------------
    ↑ update available  actions/checkout  9c091bb -> 3d3c42e (v7.0.1)
    ↑ update available  oven-sh/setup-bun  v2 -> 0c5077e (v2.2.0)
    ✓ up to date        Azure/functions-action (v1.5.6)

Check summary for <REPO1>: 10 update(s) available, 5 up to date.
```

### Project structure

| Path          | Responsibility                                                        |
| ------------- | --------------------------------------------------------------------- |
| `main.go`     | CLI entry point and command dispatch.                                 |
| `config/`     | Configuration and argument parsing, plus domain types.                |
| `handler/`    | Workflow file discovery and `uses:` dependency parsing.               |
| `operations/` | The `check` and `update` operations and shared release resolution.    |
| `exec/`       | Thin wrappers around the `gh` CLI (release tags, commit SHAs).        |
| `logger/`     | Small leveled logger used throughout the tool.                        |

## License

Licensed under the [GNU General Public License v3.0](LICENSE).
