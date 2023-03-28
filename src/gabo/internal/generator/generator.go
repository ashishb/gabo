package generator

import (
	"github.com/rs/zerolog/log"
)

type Generator struct {
	// base dir to the git repo
	dir string
	// if true, overwrite existing files
	force bool
}

func NewGenerator(dir string, force bool) Generator {
	return Generator{
		dir:   dir,
		force: force,
	}
}

func (g Generator) Generate(option Option) error {
	if g.force {
		log.Warn().Msgf("Force overwrite is on, existing files will be over-written")
	}
	switch option {
	case _LintDocker:
		return writeOrWarn(option.getOutputFileName(g.dir), option.getYamlConfig(), g.force)
	default:
		log.Fatal().Msgf("Generate for '%s' not implemented", option)
	}
	return nil
}
