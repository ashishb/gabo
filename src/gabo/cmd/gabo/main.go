package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	_modeSuggest  = "suggest"
	_modeGenerate = "generate"
)

var (
	_validModes = []string{_modeGenerate, _modeSuggest}
	_mode       = flag.String("mode", _modeSuggest,
		fmt.Sprintf("Mode to operate in: %s", _validModes))
)

func main() {
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	validateFlags()
	log.Debug().Msgf("Hello world")
}

func validateFlags() {
	if *_mode != _modeSuggest && *_mode != _modeGenerate {
		log.Fatal().Msgf("Invalid mode: %s, valid values are %s",
			*_mode, _validModes)
	}
}
