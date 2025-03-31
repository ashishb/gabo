package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashishb/gabo/src/gabo/internal/analyzer"
	"github.com/ashishb/gabo/src/gabo/internal/generator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	_options = flag.String("for", "", fmt.Sprintf("Generate GitHub Action (options: %s)",
		strings.Join(generator.GetOptionFlags(), ",")))
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
	err := validateFlags()
	if err != nil {
		log.Fatal().Msgf("%v", err.Error())
	}

	if *_version {
		fmt.Printf("gabo %s by Ashish Bhatia\nhttps://github.com/ashishb/gabo\n\n",
			strings.TrimSpace(_versionCode))
		flag.Usage()
		return
	}
	err = validateGitDir()
	if err != nil {
		log.Fatal().Msgf("%v", err.Error())
	}
	switch *_mode {
	case _modeAnalyze:
		log.Info().Msgf("Analyzing dir '%s'", *_gitDir)
		err := analyzer.Analyze(*_gitDir)
		if err != nil {
			log.Fatal().Msgf("Failed to analyze: %s", err.Error())
		}
	case _modeGenerate:
		err := generator.NewGenerator(*_gitDir, *_force).Generate(strings.Split(*_options, ","))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to generate")
		}
	}
}

// validateFlags validates flags
// This will normalize values of certain flags like _gitDir as well
func validateFlags() error {
	if *_mode != _modeAnalyze && *_mode != _modeGenerate && !*_version {
		return fmt.Errorf("invalid mode: %s, valid values are %s",
			*_mode, _validModes)
	}
	if *_force && *_mode != _modeGenerate {
		return fmt.Errorf("force overwrite is only supported in %s mode", _modeGenerate)
	}
	if *_mode == _modeGenerate {
		if _options == nil {
			return fmt.Errorf("'for' not provided in in %s mode", _modeGenerate)
		}
		options := strings.Split(*_options, ",")
		for _, option := range options {
			if !generator.IsValid(option) {
				return fmt.Errorf("'for' is not valid, valid options are one or more of %s",
					strings.Join(generator.GetOptionFlags(), ","))
			}
		}
	}
	return nil
}

// validateGitDir validates the provided dir is a git directory
func validateGitDir() error {
	if len(*_gitDir) == 0 {
		return fmt.Errorf("dir cannot be empty")
	}
	if *_gitDir == "." {
		path, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get current dir")
		}
		_gitDir = &path
	} else if strings.HasPrefix(*_gitDir, "~/") {
		tmp := strings.ReplaceAll(*_gitDir, "~", os.Getenv("HOME"))
		_gitDir = &tmp
	}
	if _, err := os.Stat(filepath.Join(*_gitDir, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("dir exists but is not a git directory: %s", *_gitDir)
	}
	return nil
}
