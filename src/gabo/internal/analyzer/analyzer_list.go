package analyzer

import "strings"

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

func contains(yamlData []string, pattern string) bool {
	for _, data := range yamlData {
		if strings.Contains(data, pattern) {
			return true
		}
	}
	return false
}
