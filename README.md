# GitHub Actions Boilerplate (`gabo`)

[![Lint YAML](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml) [![Lint Markdown](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml)
[![Lint GitHub Actions](https://github.com/ashishb/gabo/actions/workflows/lint-github-actions.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-github-actions.yaml)
[![Lint Go](https://github.com/ashishb/gabo/actions/workflows/lint-go.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-go.yaml) [![Validate Go code formatting](https://github.com/ashishb/gabo/actions/workflows/format-go.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/format-go.yaml) [![Test](https://github.com/ashishb/gabo/actions/workflows/test-go.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/test-go.yaml)

[![Check Go Releaser config for validity](https://github.com/ashishb/gabo/actions/workflows/check-goreleaser-config.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/check-goreleaser-config.yaml)

[![Release Go binaries with Go Releaser](https://github.com/ashishb/gabo/actions/workflows/release-binary.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/release-binary.yaml) [![Go report](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/ashishb/gabo/src/gabo)

**gabo** short for GitHub Actions Boilerplate is for ease-of-generation of GitHub actions boilerplate with good timeouts, path filters, and concurrency preventions. See [this blogpost](https://ashishb.net/tech/common-pitfalls-of-github-actions/) for more details on why the GitHub defaults aren't great.

The actions runs only on push/pull requests against `main` and `master` branch, by default.
Feel free to add more branches if you want to runs these checks when push/pull request against any other branches.

## Installation

```bash
$ go install github.com/ashishb/gabo/src/gabo/cmd/gabo@latest
...
```

or via homebrew ![homebrew version](https://img.shields.io/homebrew/v/gabo)

```bash
$ brew install gabo
...
```

Or run directly

```bash
# --dir flag is optional and defaults to current directory
go run github.com/ashishb/gabo/src/gabo/cmd/gabo@latest --dir=<path-to-git-dir>
```

## Usage

```bash
$ gabo --help

Usage of gabo:
  -dir string
   Path to root of git directory (default ".")
  -for string
   Generate GitHub Action (options: build-android,lint-android,translate-android,compress-images,build-docker,build-npm,build-yarn,lint-docker,format-go,lint-go,check-go-releaser,lint-html,lint-markdown,validate-openapi-schema,format-python,lint-python,lint-shell-script,lint-solidity,lint-yaml,lint-github-actions,validate-render-blueprint)
  -force
   Force overwrite existing files (in generate mode)
  -mode string
   Mode to operate in: [generate analyze] (default "analyze")
  -verbose
   Enable verbose logging
  -version
   Prints version of the binary
```

### Sample usage - analyze a repository

```bash
# Analyze current directory (it should be the root of the repository)
$ gabo

# Analyze a different dir
$ gabo --dir=~/src/repo1
```

### Sample usage - generate code

```bash
$ gabo --mode=generate --for=lint-docker
Wrote file .github/workflows/lint-docker.yaml
```

### Supported actions

- [x] build-android
- [x] build-docker
- [x] check-go-releaser
- [x] compress-images
- [x] format-go
- [x] format-python
- [x] lint-android
- [x] lint-docker
- [x] lint-go
- [x] lint-html
- [x] lint-markdown
- [x] lint-python
- [x] lint-shell-script
- [x] lint-solidity
- [x] lint-yaml
- [x] lint-github-actions
- [x] translate-android
- [x] validate-openapi-schema
- [x] validate-render-blueprint
- [ ] build-rust
- [ ] lint-rust
