package analyzer

import (
	"github.com/rs/zerolog/log"
	"strings"
)

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

func isDockerLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "hadolint")
}

func isPythonLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "pylint")
}

func isShellScriptLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "ludeeus/action-shellcheck")
}

func isSolidityLinterImplemented(yamlData []string) bool {
	// This should be made more accurate over time
	return contains(yamlData, "contractshark/inject-solhint-ci")
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
