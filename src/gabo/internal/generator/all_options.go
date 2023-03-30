package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

// Option represents a generation option
type Option string

const (
	_BuildAndroid Option = "build-android"
	_BuildDocker  Option = "build-docker"

	_FormatGo Option = "format-go"

	_LintAndroid     Option = "lint-android"
	_LintDocker      Option = "lint-docker"
	_LintGo          Option = "lint-go"
	_LintMarkdown    Option = "lint-markdown"
	_LintPython      Option = "lint-python"
	_LintShellScript Option = "lint-shell-script"
	_LintSolidity    Option = "lint-solidity"
	_LintYaml        Option = "lint-yaml"

	_TranslateAndroid Option = "translate-android"

	// _LintHtml     Option   = "lint-html"

	// TODO(ashishb): Enable these
	// _BuildGo      Option = "build-go"
	// _TestGo Option = "test-go"
	// _TestPython Option = "test-python"
)

var _options = []Option{
	_BuildAndroid,
	_BuildDocker,

	_FormatGo,

	_LintAndroid,
	_LintDocker,
	_LintGo,
	_LintMarkdown,
	_LintPython,
	_LintShellScript,
	_LintSolidity,
	_LintYaml,

	_TranslateAndroid,
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
	case _BuildAndroid:
		return getPath(repoDir, "build-android.yaml")
	case _BuildDocker:
		return getPath(repoDir, "build-docker.yaml")
	case _FormatGo:
		return getPath(repoDir, "format-go.yaml")
	case _LintAndroid:
		return getPath(repoDir, "lint-android.yaml")
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

	case _TranslateAndroid:
		return getPath(repoDir, "translate-android.yaml")
	default:
		log.Panic().Msgf("unexpected case: %s ", option)
		return ""
	}
}

func (option Option) getYamlConfig(repoDir string) (*string, error) {
	switch option {
	case _BuildAndroid:
		return generateBuildAndroidYaml(repoDir)
	case _BuildDocker:
		return generateBuildDockerYaml(repoDir)
	case _FormatGo:
		return &_formatGoYaml, nil
	case _LintAndroid:
		return &_lintAndroidYaml, nil
	case _LintDocker:
		return &_lintDockerYaml, nil
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

	case _TranslateAndroid:
		return &_translateAndroidYaml, nil

	default:
		return nil, fmt.Errorf("unexpected case: %s ", option)
	}
}

func getPath(rootDir string, fileName string) string {
	return filepath.Join(rootDir, ".github", "workflows", fileName)
}
