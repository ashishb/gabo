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
