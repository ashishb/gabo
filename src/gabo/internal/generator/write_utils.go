package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

func writeOrWarn(outFilePath string, data string, force bool) error {
	if !force && fileExists(outFilePath) {
		return fmt.Errorf("cannot write %s, file already exists", outFilePath)
	}
	if !dirExists(filepath.Dir(outFilePath)) {
		log.Debug().Msgf("Creating directory %s", filepath.Dir(outFilePath))
		err := os.MkdirAll(filepath.Dir(outFilePath), 0755)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	log.Info().Msgf("Wrote file %s", outFilePath)
	return file.Close()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if err == nil && !info.IsDir() {
		log.Warn().Msgf("Path exists but is not a directory: %s", path)
		return false
	}
	return true
}
