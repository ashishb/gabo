package analyzer

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
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
	tools["Dockerfile"] = []_Analyzer{
		{name: "Docker Linter", checker: isDockerLinterImplemented},
	}
	tools["py"] = []_Analyzer{
		{name: "Python Linter", checker: isPythonLinterImplemented},
	}
	tools["sh"] = []_Analyzer{
		{name: "Shellscript Linter", checker: isShellScriptLinterImplemented},
	}
	tools["bash"] = tools["sh"]
	tools["sol"] = []_Analyzer{
		{name: "Solidity Linter", checker: isSolidityLinterImplemented},
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
