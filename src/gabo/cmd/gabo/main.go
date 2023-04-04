package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/ashishb/gabo/src/gabo/internal/analyzer"
	"github.com/ashishb/gabo/src/gabo/internal/generator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

const (
	_modeAnalyze  = "analyze"
	_modeGenerate = "generate"
)

var (
	_verbose = flag.Bool("verbose", false, "Enable verbose logging")

	_validModes = []string{_modeGenerate, _modeAnalyze}
	_mode       = flag.String("mode", _modeAnalyze,
		fmt.Sprintf("Mode to operate in: %s", _validModes))
	_gitDir = flag.String("dir", ".", "Path to root of git directory")

	_option = flag.String("for", "", fmt.Sprintf("Generate GitHub Action (options: %s)",
		strings.Join(generator.GetOptions(), ",")))
	_force = flag.Bool("force", false,
		fmt.Sprintf("Force overwrite existing files (in %s mode)", _modeGenerate))
	_version = flag.Bool("version", false, "Prints version of the binary")
)

//go:embed version.txt
var _versionCode string

func main() {
	flag.Parse()
	originalUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("Generates GitHub Actions boilerplate\n")
		originalUsage()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if *_verbose {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	validateFlags()
	if *_version {
		fmt.Printf("gabo %s by Ashish Bhatia\nhttps://github.com/ashishb/gabo\n\n",
			strings.TrimSpace(_versionCode))
		flag.Usage()
		return
	}
	validateGitDir()

	switch *_mode {
	case _modeAnalyze:
		log.Info().Msgf("Analyzing dir '%s'", *_gitDir)
		analyzer.Analyze(*_gitDir)
	case _modeGenerate:
		err := generator.NewGenerator(*_gitDir, *_force).Generate(
			generator.Option(*_option))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to generate")
		}
	}
}

// This will normalize values of certain flags like _gitDir as well
func validateFlags() {
	if *_mode != _modeAnalyze && *_mode != _modeGenerate && !*_version {
		log.Fatal().Msgf("Invalid mode: %s, valid values are %s",
			*_mode, _validModes)
		return
	}
	if *_force && *_mode != _modeGenerate {
		log.Fatal().Msgf("force overwrite is only supported in %s mode", _modeGenerate)
		return
	}
	if *_mode == _modeGenerate {
		if _option == nil {
			log.Fatal().Msgf("'for' not provided in in %s mode", _modeGenerate)
			return
		}
		if !generator.IsValid(*_option) {
			log.Fatal().Msgf("'for' is not valid, valid options are %s",
				strings.Join(generator.GetOptions(), ","))
			return
		}
	}
}

func validateGitDir() bool {
	if len(*_gitDir) == 0 {
		log.Fatal().Msgf("dir cannot be empty")
		return false
	}
	if *_gitDir == "." {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to get current dir")
		}
		_gitDir = &path
	} else if strings.HasPrefix(*_gitDir, "~/") {
		tmp := strings.ReplaceAll(*_gitDir, "~", os.Getenv("HOME"))
		_gitDir = &tmp
	}
	if _, err := os.Stat(filepath.Join(*_gitDir, ".git")); os.IsNotExist(err) {
		log.Fatal().Msgf("dir exists but is not a git directory: %s", *_gitDir)
		return false
	}
	return true
}
