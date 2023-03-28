package generator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func writeOrWarn(outFilePath string, data string, force bool) error {
	if !force && fileExists(outFilePath) {
		return fmt.Errorf("cannot write %s, file already exists", outFilePath)
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
