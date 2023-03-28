package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

// Option represents a generation option
type Option string

const (
	_LintDocker      Option = "lint-docker"
	_LintGo          Option = "lint-go"
	_LintMarkdown    Option = "lint-markdown"
	_LintPython      Option = "lint-python"
	_LintShellScript Option = "lint-shell-script"
	_LintSolidity    Option = "lint-solidity"
	_LintYaml        Option = "lint-yaml"

	_BuildAndroid Option = "build-android"
	_BuildDocker  Option = "build-docker"
	_BuildGo      Option = "build-go"
	_BuildPython  Option = "build-python"

	// Make this code correct first
	// _LintAndroid  Option   = "lint-android"
	// _LintHtml     Option   = "lint-html"

	// TODO(ashishb): Enable these
	// _AutoTranslateAndroid Option = "auto-translate-android"
	// _TestGo Option = "test-go"
	// _TestPython Option = "test-python"
)

var _options = []Option{
	_LintDocker,
	_LintGo,
	_LintMarkdown,
	_LintPython,
	_LintShellScript,
	_LintSolidity,
	_LintYaml,

	_BuildAndroid,
	_BuildDocker,
	_BuildGo,
	_BuildPython,
}

func GetOptions() []string {
	result := make([]string, 0, len(_options))
	for _, option := range _options {
		result = append(result, string(option))
	}
	return result
}

func IsValid(val string) bool {
	for _, option := range _options {
		if val == string(option) {
			return true
		}
	}
	return false
}

// repoDir is root dir of the repository
func (option Option) getOutputFileName(repoDir string) string {
	switch option {
	case _LintDocker:
		return getPath(repoDir, "lint-docker.yaml")
	case _LintGo:
		return getPath(repoDir, "lint-go.yaml")
	case _LintMarkdown:
		return getPath(repoDir, "lint-markdown.yaml")
	case _LintPython:
		return getPath(repoDir, "lint-python.yaml")
	case _LintShellScript:
		return getPath(repoDir, "lint-shell-script.yaml")
	case _LintSolidity:
		return getPath(repoDir, "lint-solidity.yaml")
	case _LintYaml:
		return getPath(repoDir, "lint-yaml.yaml")

	case _BuildDocker:
		return getPath(repoDir, "build-docker.yaml")
	case _BuildGo:
		return getPath(repoDir, "build-go.yaml")
	case _BuildPython:
		return getPath(repoDir, "build-python.yaml")
	default:
		log.Panic().Msgf("unexpected case: %s ", option)
		return ""
	}
}

func (option Option) getYamlConfig(repoDir string) (*string, error) {
	switch option {
	case _LintDocker:
		return &_lintDockerYaml, nil
	// TODO(ashishb): This does not handle monorepo case just yet
	// fix that
	case _LintGo:
		return generateGoLintYaml(repoDir)
	case _LintMarkdown:
		return &_lintMarkdownYaml, nil
	case _LintPython:
		return &_lintPythonYaml, nil
	case _LintShellScript:
		return &_lintShellScriptYaml, nil
	case _LintSolidity:
		return &_lintSolidityYaml, nil
	case _LintYaml:
		return &_lintYamlYaml, nil

	case _BuildDocker:
		fallthrough
	case _BuildGo:
		fallthrough
	case _BuildPython:
		fallthrough
	default:
		return nil, fmt.Errorf("unexpected case: %s ", option)
	}
}

func getPath(rootDir string, fileName string) string {
	return filepath.Join(rootDir, ".github", "workflows", fileName)
}
