package analyzer

import (
	"fmt"
	"github.com/ashishb/gabo/src/gabo/internal/generator"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type _Analyzer struct {
	name    string
	option  generator.Option
	checker func(yamlStrings []string) bool
}

func (a _Analyzer) generateCommand(rootDir string) interface{} {
	if strings.Contains(rootDir, " ") {
		// escape whitespace in rootdir
		rootDir = fmt.Sprintf("'%s'", rootDir)
	}
	return fmt.Sprintf("%s --mode=generate --for=%s --dir=%s", os.Args[0], a.option, rootDir)
}

func (a _Analyzer) isApplicable(rootDir string) bool {
	return a.option.IsApplicable(rootDir)
}

func Analyze(rootDir string) {
	workflowsDir := filepath.Join(rootDir, ".github", "workflows")
	yamlStrings, err := getYamlData(workflowsDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}
	analyzers := []_Analyzer{
		{"Android Linter", generator.LintAndroid, isAndroidLinterImplemented},
		{"Android Auto-translate", generator.TranslateAndroid, isAndroidAutoTranslatorImplemented},
		{"Docker Linter", generator.LintDocker, isDockerLinterImplemented},
		{"Markdown Linter", generator.LintMarkdown, isMarkdownLinterImplemented},
		{"Go Linter", generator.LintGo, isGoLinterImplemented},
		{"Go Formatter", generator.FormatGo, isGoFormatterImplemented},
		{"OpenAPI Schema Validator", generator.ValidateOpenApiSchema, isOpenApiSchemaValidatorImplemented},
		{"Python Linter", generator.LintPython, isPythonLinterImplemented},
		{"Shellscript Linter", generator.LintShellScript, isShellScriptLinterImplemented},
		{"Solidity Linter", generator.LintSolidity, isSolidityLinterImplemented},
		{"YAML Linter", generator.LintYaml, isYamlLinterImplemented},
	}

	suggestedCount := 0
	for _, analyzer := range analyzers {
		if !analyzer.isApplicable(rootDir) {
			log.Trace().Msgf("Not applicable %s", analyzer.name)
			continue
		}
		if analyzer.checker(yamlStrings) {
			log.Info().Msgf("✅ %s is present", analyzer.name)
		} else {
			log.Warn().Msgf("❌ %s is missing, (\"%s\")",
				analyzer.name, analyzer.generateCommand(rootDir))
			suggestedCount = 0
		}
	}
	if suggestedCount == 0 {
		log.Info().Msg("No Actions changes")
	}
}

func getYamlData(dir string) ([]string, error) {
	// Glob pattern "*.y?(a)ml" is not supported by Go
	globPattern1 := filepath.Join(dir, "*.yaml")
	globPattern2 := filepath.Join(dir, "*.yml")
	log.Trace().Msgf("Glob patterns are %s and %s", globPattern1, globPattern2)
	matches1, err1 := filepath.Glob(globPattern1)
	matches2, err2 := filepath.Glob(globPattern2)
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		return nil, nil
	}
	matches := make([]string, 0)
	if len(matches1) > 0 {
		matches = append(matches, matches1...)
	}
	if len(matches2) > 0 {
		matches = append(matches, matches2...)
	}
	data := make([]string, 0, len(matches))
	for _, filePath := range matches {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		tmp, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		log.Trace().Msgf("File %s", filePath)
		data = append(data, string(tmp))
	}
	return data, nil
}
