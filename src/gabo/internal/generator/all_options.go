package generator

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

type Option interface {
	// E.g. "Markdown Linter"
	Name() string
	// E.g. lint-markdown
	FlagName() string
	IsApplicable(repoDir string) bool
	IsImplemented(yamlData []string) bool
	GetYamlConfig(repoDir string) (*string, error)
	GetOutputFileName(repoDir string) string
}

func GetOptions() []Option {
	// _BuildGo      Option = "build-go"
	// _TestGo Option = "test-go"
	// _TestPython Option = "test-python"

	return []Option{
		_Option{"Android Builder", "build-android",
			newFileMatcher("AndroidManifest.xml"),
			newPatternMatcher("gradlew build"),
			newGenerator2(generateBuildAndroidYaml), "build-android.yaml"},
		_Option{"Android Linter", "lint-android",
			newFileMatcher("AndroidManifest.xml"),
			newPatternMatcher("gradlew lint"),
			newGenerator(_lintAndroidYaml), "lint-android.yaml"},
		_Option{"Android Auto Translator", "translate-android",
			newFileMatcher("AndroidManifest.xml"),
			newPatternMatcher("ashishb/android-auto-translate"),
			newGenerator(_translateAndroidYaml),
			"translate-android.yaml"},

		_Option{"Docker Builder", "build-docker",
			newFileMatcher("Dockerfile"),
			newPatternMatcher("docker build "),
			newGenerator2(generateBuildDockerYaml), "build-docker.yaml"},
		_Option{"Docker Linter", "lint-docker", newFileMatcher("Dockerfile"),
			newPatternMatcher("hadolint "),
			newGenerator(_lintDockerYaml), "lint-docker.yaml"},

		_Option{"Go Formatter", "format-go", newFileMatcher("*.go"),
			newPatternMatcher("gofmt -l", "gofumpt "),
			newGenerator(_formatGoYaml), "format-go.yaml"},
		_Option{"Go Linter", "lint-go", newFileMatcher("*.go"),
			newPatternMatcher("golangci-lint run "),
			newGenerator2(generateGoLintYaml), "lint-go.yaml"},

		_Option{"HTML Linter", "lint-html", newFileMatcher("*.html", "*.htm"),
			newPatternMatcher("htmlhint "), newGenerator(_lintHtmlYaml),
			"lint-html.yaml"},

		_Option{"Markdown Linter", "lint-markdown", newFileMatcher("*.md"),
			newPatternMatcher("mdl "),
			newGenerator(_lintMarkdownYaml), "lint-markdown.yaml"},
		_Option{"OpenAPI Schema Validator", "validate-openapi-schema",
			newFileMatcher("openapi.json", "openapi.yaml", "openapi.yml"),
			newPatternMatcher("swagger-cli validate "),
			newGenerator2(generateOpenAPISchemaValidator),
			"validate-openapi-schema.yaml"},
		_Option{"Python Linter", "lint-python", newFileMatcher("*.py"),
			newPatternMatcher("pylint "),
			newGenerator(_lintPythonYaml), "lint-python.yaml"},
		_Option{"Shell Script Linter", "lint-shell-script", newFileMatcher("*.sh", "*.bash"),
			newPatternMatcher("shellcheck "),
			newGenerator(_lintShellScriptYaml), "lint-shell-script.yaml"},
		_Option{"Solidity Linter", "lint-solidity", newFileMatcher("*.sol"),
			newPatternMatcher("solhint "),
			newGenerator(_lintSolidityYaml), "lint-solidity.yaml"},
		_Option{"YAML Linter", "lint-yaml", newFileMatcher("*.yml", "*.yaml"),
			newPatternMatcher("yamllint "),
			newGenerator(_lintYamlYaml), "lint-yaml.yaml"},
	}
}

type _FileMatcher interface {
	Matches(repoDir string) bool
}

type _FileMatcherImpl struct {
	patterns []string
}

func (f _FileMatcherImpl) Matches(rootDir string) bool {
	for _, pattern := range f.patterns {
		if hasFile(rootDir, pattern) {
			return true
		}
	}
	return false
}

func newFileMatcher(patterns ...string) _FileMatcher {
	return _FileMatcherImpl{patterns}
}

type _PatternMatcher interface {
	Matches(yamlData []string) bool
}

type _PatternMatcherImpl struct {
	patterns []string
}

func (f _PatternMatcherImpl) Matches(yamlData []string) bool {
	for _, pattern := range f.patterns {
		if contains(yamlData, pattern) {
			return true
		}
	}
	return false
}

func newPatternMatcher(patterns ...string) _PatternMatcher {
	return _PatternMatcherImpl{patterns}
}

type _Generator interface {
	Generate(repoDir string) (*string, error)
}

type _GeneratorStringImpl struct {
	yamlStr string
}

func (g _GeneratorStringImpl) Generate(_ string) (*string, error) {
	return &g.yamlStr, nil
}

type _GeneratorFuncImpl struct {
	f func(repoDir string) (*string, error)
}

func (g _GeneratorFuncImpl) Generate(repoDir string) (*string, error) {
	return g.f(repoDir)
}

func newGenerator(yamlStr string) _Generator {
	return _GeneratorStringImpl{yamlStr}
}

func newGenerator2(f func(repoDir string) (*string, error)) _Generator {
	return _GeneratorFuncImpl{f: f}
}

func GetOptionFlags() []string {
	result := make([]string, 0)
	for _, option := range GetOptions() {
		result = append(result, option.FlagName())
	}
	return result
}

func IsValid(val string) bool {
	for _, option := range GetOptionFlags() {
		if val == option {
			return true
		}
	}
	return false
}

type _Option struct {
	name                        string
	flagName                    string
	filePatternToCheck          _FileMatcher
	isImplementedPatternMatcher _PatternMatcher
	yamlConfigGenerator         _Generator
	outputFileName              string
}

func (o _Option) Name() string {
	return o.name
}

func (o _Option) FlagName() string {
	return o.flagName
}

func (o _Option) IsApplicable(repoDir string) bool {
	return o.filePatternToCheck.Matches(repoDir)
}

func (o _Option) IsImplemented(yamlData []string) bool {
	return o.isImplementedPatternMatcher.Matches(yamlData)
}

func (o _Option) GetYamlConfig(repoDir string) (*string, error) {
	return o.yamlConfigGenerator.Generate(repoDir)
}

func (o _Option) GetOutputFileName(repoDir string) string {
	return getPath(repoDir, o.outputFileName)
}

func getPath(rootDir string, fileName string) string {
	return filepath.Join(rootDir, ".github", "workflows", fileName)
}

// Note: Go does not support "**" glob pattern
func hasFile(rootDir string, globPattern string) bool {
	found := false
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		matched, err := filepath.Match(globPattern, info.Name())
		if err != nil {
			return err
		}
		if matched {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		log.Panic().Err(err).Msgf("glob failed: '%s'", globPattern)
	}
	log.Trace().Msgf("hasFile(%s, %s) = %v", globPattern, rootDir, found)
	return found
}

func contains(yamlData []string, pattern string) bool {
	log.Trace().Msgf("Looking for %s", pattern)
	for _, data := range yamlData {
		if strings.Contains(data, pattern) {
			return true
		}
	}
	return false
}
