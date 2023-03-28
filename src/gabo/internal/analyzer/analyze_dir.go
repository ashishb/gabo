package analyzer

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type _Analyzer struct {
	name    string
	checker func(yamlStrings []string) bool
}

func Analyze(rootDir string) {
	workflowsDir := filepath.Join(rootDir, ".github", "workflows")
	yamlStrings, err := getYamlData(workflowsDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}
	extToFreqMap := make(map[string]int)
	// For files like "Dockerfile", extension would be the full
	// file name
	err = getExtToFileCountMap(extToFreqMap, rootDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}
	for k, v := range extToFreqMap {
		log.Debug().Msgf("%s\t%d", k, v)
	}
	// ext -> array of analyzer valid for that
	tools := make(map[string][]_Analyzer)
	tools["yaml"] = []_Analyzer{{name: "YAML Linter", checker: isYamlLinterImplemented}}
	tools["md"] = []_Analyzer{{name: "Markdown Linter", checker: isMarkdownLinterImplemented}}
	tools["go"] = []_Analyzer{
		{name: "Go Linter", checker: isGoLinterImplemented},
		{name: "Go Formatter", checker: isGoFormatterImplemented},
	}

	for ext, analyzers := range tools {
		if extToFreqMap[ext] <= 0 {
			continue
		}
		for _, analyzer := range analyzers {
			if analyzer.checker(yamlStrings) {
				log.Info().Msgf("✅ %s is present", analyzer.name)
			} else {
				log.Warn().Msgf("❌ %s is missing", analyzer.name)
			}
		}
	}
}

func getYamlData(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "**.y[a]ml"))
	if os.IsNotExist(err) {
		return nil, nil
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
		data = append(data, string(tmp))
	}
	return data, nil
}

func isYamlLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "uses: ibiqlik/action-yamllint")
}

func isMarkdownLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "mdl ")
}

func isGoLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "uses: golangci/golangci-lint-action")
}

func isGoFormatterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "gofmt -l")
}

func contains(yamlData []string, pattern string) bool {
	for _, data := range yamlData {
		if strings.Contains(data, pattern) {
			return true
		}
	}
	return false
}
