package main

import (
	"flag"
	"fmt"
	"github.com/ashishb/gabo/src/gabo/internal/analyzer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

const (
	_modeSuggest  = "suggest"
	_modeGenerate = "generate"
)

var (
	_validModes = []string{_modeGenerate, _modeSuggest}
	_mode       = flag.String("mode", _modeSuggest,
		fmt.Sprintf("Mode to operate in: %s", _validModes))
	_gitDir  = flag.String("dir", ".", "Path to root of git directory")
	_verbose = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if *_verbose {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	validateFlags()
	switch *_mode {
	case _modeSuggest:
		log.Info().Msgf("Analyzing dir '%s'", *_gitDir)
		analyzer.Analyze(*_gitDir)
	case _modeGenerate:
		log.Fatal().Msgf("Mode not implemented yet: %s", _modeGenerate)
	}
}

// This will normalize values of certain flags like _gitDir as well
func validateFlags() {
	if *_mode != _modeSuggest && *_mode != _modeGenerate {
		log.Fatal().Msgf("Invalid mode: %s, valid values are %s",
			*_mode, _validModes)
		return
	}
	if len(*_gitDir) == 0 {
		log.Fatal().Msgf("dir cannot be empty")
		return
	}
	if *_gitDir == "." {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to get current dir")
		}
		_gitDir = &path
	}
	if _, err := os.Stat(filepath.Join(*_gitDir, ".git")); os.IsNotExist(err) {
		log.Fatal().Msgf("dir exists but is not a git directory: %s", *_gitDir)
		return
	}
}
