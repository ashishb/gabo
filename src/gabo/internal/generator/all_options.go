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

	FormatGo Option = "format-go"

	LintAndroid     Option = "lint-android"
	LintDocker      Option = "lint-docker"
	LintGo          Option = "lint-go"
	LintMarkdown    Option = "lint-markdown"
	LintPython      Option = "lint-python"
	LintShellScript Option = "lint-shell-script"
	LintSolidity    Option = "lint-solidity"
	LintYaml        Option = "lint-yaml"

	TranslateAndroid      Option = "translate-android"
	ValidateOpenApiSchema Option = "validate-openapi"

	// _LintHtml     Option   = "lint-html"

	// TODO(ashishb): Enable these
	// _BuildGo      Option = "build-go"
	// _TestGo Option = "test-go"
	// _TestPython Option = "test-python"
)

var _options = []Option{
	_BuildAndroid,
	_BuildDocker,

	FormatGo,

	LintAndroid,
	LintDocker,
	LintGo,
	LintMarkdown,
	LintPython,
	LintShellScript,
	LintSolidity,
	LintYaml,

	TranslateAndroid,
	ValidateOpenApiSchema,
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
	case FormatGo:
		return getPath(repoDir, "format-go.yaml")
	case LintAndroid:
		return getPath(repoDir, "lint-android.yaml")
	case LintDocker:
		return getPath(repoDir, "lint-docker.yaml")
	case LintGo:
		return getPath(repoDir, "lint-go.yaml")
	case LintMarkdown:
		return getPath(repoDir, "lint-markdown.yaml")
	case LintPython:
		return getPath(repoDir, "lint-python.yaml")
	case LintShellScript:
		return getPath(repoDir, "lint-shell-script.yaml")
	case LintSolidity:
		return getPath(repoDir, "lint-solidity.yaml")
	case LintYaml:
		return getPath(repoDir, "lint-yaml.yaml")

	case TranslateAndroid:
		return getPath(repoDir, "translate-android.yaml")
	case ValidateOpenApiSchema:
		return getPath(repoDir, "validate-openapi-schema.yaml")
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
	case FormatGo:
		return &_formatGoYaml, nil
	case LintAndroid:
		return &_lintAndroidYaml, nil
	case LintDocker:
		return &_lintDockerYaml, nil
	case LintGo:
		return generateGoLintYaml(repoDir)
	case LintMarkdown:
		return &_lintMarkdownYaml, nil
	case LintPython:
		return &_lintPythonYaml, nil
	case LintShellScript:
		return &_lintShellScriptYaml, nil
	case LintSolidity:
		return &_lintSolidityYaml, nil
	case LintYaml:
		return &_lintYamlYaml, nil
	case TranslateAndroid:
		return &_translateAndroidYaml, nil
	case ValidateOpenApiSchema:
		return generateOpenAPISchemaValidator(repoDir)
	default:
		return nil, fmt.Errorf("unexpected case: %s ", option)
	}
}

func (option Option) IsApplicable(dir string) bool {
	switch option {
	case _BuildAndroid:
		return hasFile("**/build.gradle", dir)
	case _BuildDocker:
		return hasFile("Dockerfile", dir)
	case FormatGo:
		return hasFile("**/*.go", dir)
	case LintAndroid:
		return hasFile("**/build.gradle", dir)
	case LintDocker:
		return hasFile("Dockerfile", dir)
	case LintGo:
		return hasFile("**/*.go", dir)
	case LintMarkdown:
		return hasFile("**/*.md", dir)
	case LintPython:
		return hasFile("**/*.py", dir)
	case LintShellScript:
		return hasFile("**/*.sh", dir) || hasFile("**/*.bash", dir)
	case LintSolidity:
		return hasFile("**/*.sol", dir)
	case LintYaml:
		return hasFile("**/*.yaml", dir) || hasFile("**/*.yml", dir)
	case TranslateAndroid:
		return hasFile("**/build.gradle", dir)
	case ValidateOpenApiSchema:
		return hasFile("openapi.json", dir) ||
			hasFile("openapi.yaml", dir) ||
			hasFile("openapi.yml", dir)
	default:
		log.Panic().Msgf("unexpected case: %s ", option)
		return false
	}
}

func getPath(rootDir string, fileName string) string {
	return filepath.Join(rootDir, ".github", "workflows", fileName)
}

func hasFile(globPattern string, rootDir string) bool {
	matches, err := filepath.Glob(rootDir + "/" + globPattern)
	if err != nil {
		log.Panic().Err(err).Msgf("glob failed: '%s'", globPattern)
	}
	return len(matches) > 0
}
