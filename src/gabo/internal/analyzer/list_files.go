package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

func listAllFiles(result map[string]int, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		// Ignore
		if entry.IsDir() && (entry.Name() == ".git" || entry.Name() == "node_modules" ||
			entry.Name() == ".idea") {
			continue
		}
		if entry.IsDir() {
			err = listAllFiles(result, filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
			continue
		}
		fields := strings.Split(entry.Name(), ".")
		ext := fields[len(fields)-1]
		// Normalize extension
		if ext == "htm" {
			ext = "html"
		}
		if ext == "yml" {
			ext = "yaml"
		}
		result[ext] = result[ext] + 1
	}
	return nil
}
