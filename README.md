# GitHub Actions Boilerplate (`gabo`)

[![Lint YAML](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-yaml.yaml) [![Lint Markdown](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml/badge.svg)](https://github.com/ashishb/gabo/actions/workflows/lint-markdown.yaml)

**gabo** short for GitHub Actions Boilerplate is for ease-of-generation of GitHub actions boilerplate with good timeouts, path filters, and concurrency preventions. See [this blogpost](https://ashishb.net/tech/common-pitfalls-of-github-actions/) for more details on why the GitHub defaults aren't great.

The actions runs only on push/pull requests against `main` and `master` branch, by default.
Feel free to add more branches if you want to runs these checks when push/pull request against any other branches.
