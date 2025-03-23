package generator

import (
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	_AndroidManifestFile = "**/AndroidManifest.xml"
	_dockerFile          = "**/Dockerfile"
	_goFile              = "**/*.go"
	_markdownFile        = "**/*.md"
	_pythonFile          = "**/*.py"
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
		_Option{
			"Android Builder", "build-android",
			newFileMatcher(_AndroidManifestFile),
			newPatternMatcher("gradlew build"),
			newGenerator2(generateBuildAndroidYaml), "build-android.yaml",
		},
		_Option{
			"Android Linter", "lint-android",
			newFileMatcher(_AndroidManifestFile),
			newPatternMatcher("gradlew lint"),
			newGenerator(_lintAndroidYaml), "lint-android.yaml",
		},
		_Option{
			"Android Auto Translator", "translate-android",
			newFileMatcher(_AndroidManifestFile),
			newPatternMatcher("ashishb/android-auto-translate"),
			newGenerator(_translateAndroidYaml),
			"translate-android.yaml",
		},
		_Option{
			"Compress Images", "compress-images",
			newFileMatcher("**/*.jpg", "**/*.jpeg", "**/*.png", "**/*.webp"),
			newPatternMatcher("calibreapp/image-actions"),
			newGenerator(_comressImageYaml), "compress-images.yaml",
		},
		_Option{
			"Docker Builder", "build-docker",
			newFileMatcher(_dockerFile),
			newPatternMatcher("docker build ", "docker buildx"),
			newGenerator2(generateBuildDockerYaml), "build-docker.yaml",
		},
		_Option{
			"NPM Builder", "build-npm", newFileMatcher("**/package-lock.json"),
			newPatternMatcher("npm install "),
			newGenerator2(generateBuildNpmYaml), "build-npm.yaml",
		},
		_Option{
			"Yarn Builder", "build-yarn", newFileMatcher("**/yarn.lock"),
			newPatternMatcher("yarn build"),
			newGenerator2(generateBuildYarnYaml), "build-yarn.yaml",
		},
		_Option{
			"Docker Linter", "lint-docker", newFileMatcher(_dockerFile),
			newPatternMatcher("hadolint "),
			newGenerator(_lintDockerYaml), "lint-docker.yaml",
		},

		_Option{
			"Go Formatter", "format-go", newFileMatcher(_goFile),
			newPatternMatcher("gofmt -l", "go fmt", "gofumpt "),
			newGenerator(_formatGoYaml), "format-go.yaml",
		},
		_Option{
			"Go Linter", "lint-go", newFileMatcher(_goFile),
			newPatternMatcher("golangci-lint"),
			newGenerator2(generateGoLintYaml), "lint-go.yaml",
		},
		_Option{
			"Go Releaser Config Checker", "check-go-releaser",
			newFileMatcher(_goReleaserConfigFiles...),
			newPatternMatcher("goreleaser check "),
			newGenerator2(generateGoReleaserConfigCheckerYaml), "check-goreleaser-config.yaml",
		},

		_Option{
			"HTML Linter", "lint-html", newFileMatcher("**/*.html", "**/*.htm"),
			newPatternMatcher("htmlhint "), newGenerator(_lintHtmlYaml),
			"lint-html.yaml",
		},

		_Option{
			"Markdown Linter", "lint-markdown", newFileMatcher(_markdownFile),
			newPatternMatcher("mdl "),
			newGenerator(_lintMarkdownYaml), "lint-markdown.yaml",
		},
		_Option{
			"OpenAPI Schema Validator", "validate-openapi-schema",
			newFileMatcher(_openAPIFileList...),
			newPatternMatcher("mpetrunic/swagger-cli-action", "dshanley/vacuum"),
			newGenerator2(generateOpenAPISchemaValidator),
			"validate-openapi-schema.yaml",
		},
		_Option{
			"Python Formatter", "format-python", newFileMatcher(_pythonFile),
			newPatternMatcher("black "),
			newGenerator(_formatPythonYaml), "format-python.yaml",
		},
		_Option{
			"Python Linter", "lint-python", newFileMatcher(_pythonFile),
			newPatternMatcher("pylint ", "ruff "),
			newGenerator(_lintPythonYaml), "lint-python.yaml",
		},
		_Option{
			"Shell Script Linter", "lint-shell-script", newFileMatcher("**/*.sh", "**/*.bash"),
			newPatternMatcher("shellcheck ", "action-shellcheck"),
			newGenerator(_lintShellScriptYaml), "lint-shell-script.yaml",
		},
		_Option{
			"Solidity Linter", "lint-solidity", newFileMatcher("**/*.sol"),
			newPatternMatcher("solhint "),
			newGenerator(_lintSolidityYaml), "lint-solidity.yaml",
		},
		_Option{
			"YAML Linter", "lint-yaml", newFileMatcher("**/*.yml", "**/*.yaml"),
			newPatternMatcher("ibiqlik/action-yamllint@"),
			newGenerator(_lintYamlYaml), "lint-yaml.yaml",
		},
		_Option{
			"GitHub Actions Linter", "lint-github-actions",
			newFileMatcher(".github/workflows/*.yml", ".github/workflows/*.yaml"),
			newPatternMatcher( /*"actionlint",*/ "zizmor"),
			newGenerator(_lintGitHubActionsYaml), "lint-github-actions.yaml",
		},
		_Option{
			"Render.com blueprint Validator", "validate-render-blueprint", newFileMatcher("render.yml", "render.yaml"),
			newPatternMatcher("GrantBirki/json-yaml-validate"),
			newGenerator(_validateRenderBlueprintYaml), "validate-render-blueprint.yaml",
		},
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
	log.Debug().
		Str("rootDir", rootDir).
		Str("globPattern", globPattern).
		Msg("Glob pattern")
	found := false
	err := doublestar.GlobWalk(os.DirFS(rootDir), globPattern, func(path string, d fs.DirEntry) error {
		found = true
		log.Trace().Msgf("hasFile(%s, %s) = %s", globPattern, rootDir, path)
		return doublestar.SkipDir
	})
	if err != nil {
		log.Error().
			Err(err).Msgf("glob failed: '%s'", globPattern)
	}
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
