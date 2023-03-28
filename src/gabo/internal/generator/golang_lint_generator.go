package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

const _goLangLintTask = `

      - name: Run golangci-lint on %s
        uses: golangci/golangci-lint-action@v3
        with:
          # Directory containing go.mod file
          working-directory: %s
`

func generateGoLintYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "go.mod")
	if err != nil {
		return nil, err
	}
	str := _lintGoYaml
	for _, dir := range dirs {
		str += fmt.Sprintf(_goLangLintTask, dir, dir)
	}
	return &str, nil
}

// Returns "relative" path of dirs inside baseDir containing "fileName"
// Error is returned instead of the empty array
func getDirsContaining(dir string, fileName string) ([]string, error) {
	absPaths, err := getDirsContaining2(dir, fileName)
	if err != nil {
		return nil, err
	}
	if len(absPaths) == 0 {
		return nil, fmt.Errorf("no dir containing %s", fileName)
	}
	return getRelativePaths(dir, absPaths), nil
}

func getRelativePaths(dir string, absPaths []string) []string {
	relPaths := make([]string, 0, len(absPaths))
	for _, absPath := range absPaths {
		relPath := strings.TrimPrefix(strings.TrimPrefix(absPath, dir), string(os.PathSeparator))
		if len(relPath) == 0 {
			relPath = "."
		}
		relPaths = append(relPaths, relPath)
	}
	return relPaths
}

// Returns "absolute" path of dirs inside baseDir containing "fileName"
// Error is returned instead of the empty array
func getDirsContaining2(dir string, fileName string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0)

	for _, entry := range entries {
		// Ignore
		if entry.IsDir() && (entry.Name() == ".git" || entry.Name() == "node_modules" ||
			entry.Name() == ".idea") {
			continue
		}
		if entry.IsDir() {
			tmp, err := getDirsContaining2(filepath.Join(dir, entry.Name()), fileName)
			if err != nil {
				return nil, err
			}
			results = append(results, tmp...)
			continue
		}
		if strings.EqualFold(entry.Name(), fileName) {
			log.Debug().Msgf("Looking at the file %s", entry.Name())
			results = append(results, dir)
		}
	}
	return results, nil
}
