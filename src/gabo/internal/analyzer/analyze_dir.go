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

func generateCommand(flagName string, rootDir string) string {
	if strings.Contains(rootDir, " ") {
		// escape whitespace in rootdir
		rootDir = fmt.Sprintf("'%s'", rootDir)
	}
	return fmt.Sprintf("%s --mode=generate --for=%s --dir=%s", os.Args[0], flagName, rootDir)
}

func Analyze(rootDir string) {
	workflowsDir := filepath.Join(rootDir, ".github", "workflows")
	yamlStrings, err := getYamlData(workflowsDir)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err.Error())
	}

	missingAnalyzers := make([]string, 0)
	for _, analyzer := range generator.GetOptions() {
		if !analyzer.IsApplicable(rootDir) {
			log.Trace().Msgf("Not applicable %s", analyzer.Name())
			continue
		}
		if analyzer.IsImplemented(yamlStrings) {
			log.Info().Msgf("✅ %s is present", analyzer.Name())
		} else {
			log.Warn().Msgf("❌ %s is missing, generate via \"%s\"",
				analyzer.Name(), generateCommand(analyzer.FlagName(), rootDir))
			missingAnalyzers = append(missingAnalyzers, analyzer.FlagName())
		}
	}
	if len(missingAnalyzers) == 0 {
		log.Info().Msg("No changes required")
		return
	}
	log.Info().Msgf("Run the following command to generate "+
		"all the suggested GitHub Actions:\n%s", generateCommand(
		strings.Join(missingAnalyzers, ","), rootDir))
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
