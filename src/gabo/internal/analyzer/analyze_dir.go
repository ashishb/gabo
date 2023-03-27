package analyzer

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Analyze(rootDir string) {
	workflowsDir := filepath.Join(rootDir, ".github", "workflows")
	yamlStrings, err := getYamlData(workflowsDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}
	extToFreqMap := make(map[string]int)
	err = listAllFiles(extToFreqMap, rootDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}
	for k, v := range extToFreqMap {
		log.Debug().Msgf("%s\t%d", k, v)
	}
	if extToFreqMap["yaml"] > 0 {
		yamlLinter := isYamlLinterImplemented(yamlStrings)
		if !yamlLinter {
			log.Warn().Msgf("❌ YAML Linter is missing")
		} else {
			log.Info().Msgf("✅ YAML Linter is present")
		}
	}
	if extToFreqMap["md"] > 0 {
		yamlLinter := isMarkdownLinterImplemented(yamlStrings)
		if !yamlLinter {
			log.Warn().Msgf("❌ Markdown Linter is missing")
		} else {
			log.Info().Msgf("✅ Markdown Linter is present")
		}
	}
	if extToFreqMap["go"] > 0 {
		goLinter := isGoLinterImplemented(yamlStrings)
		if !goLinter {
			log.Warn().Msgf("❌ Go Linter is missing")
		} else {
			log.Info().Msgf("✅ Go Linter is present")
		}
		goFormatter := isGoFormatterImplemented(yamlStrings)
		if !goFormatter {
			log.Warn().Msgf("❌ Go Formatter is missing")
		} else {
			log.Info().Msgf("✅ Go Formatter is present")
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
