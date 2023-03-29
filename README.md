# GitHub Actions Boilerplate (`gabo`)

[![Lint YAML](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml) [![Lint Markdown](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml)
[![Lint Go](https://github.com/ashishb/gabo/actions/workflows/lint-go.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-go.yaml) [![Validate Go code formatting](https://github.com/ashishb/gabo/actions/workflows/format-go.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/format-go.yaml)

[![Check Go Releaser config for validity](https://github.com/ashishb/gabo/actions/workflows/check-goreleaser-config.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/check-goreleaser-config.yaml)

[![Release Go binaries with Go Releaser](https://github.com/ashishb/gabo/actions/workflows/release-binary.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/release-binary.yaml)

**gabo** short for GitHub Actions Boilerplate is for ease-of-generation of GitHub actions boilerplate with good timeouts, path filters, and concurrency preventions. See [this blogpost](https://ashishb.net/tech/common-pitfalls-of-github-actions/) for more details on why the GitHub defaults aren't great.

The actions runs only on push/pull requests against `main` and `master` branch, by default.
Feel free to add more branches if you want to runs these checks when push/pull request against any other branches.

## Installation

```bash
go install github.com/ashishb/gabo@latest
```

## Usage

```bash
$ gabo --help
  -dir string
      Path to root of git directory (default ".")
  -for string
      Generate GitHub Action (options: translate-android,build-android,build-docker,lint-android,lint-docker,lint-go,lint-markdown,lint-python,lint-shell-script,lint-solidity,lint-yaml)
  -force
      Force overwrite existing files (in generate mode)
  -mode string
      Mode to operate in: [generate analyze] (default "analyze")
  -verbose
      Enable verbose logging
```
