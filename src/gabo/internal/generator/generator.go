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
	str, err := option.getYamlConfig(g.dir)
	if err != nil {
		return err
	}
	return writeOrWarn(option.getOutputFileName(g.dir), *str, g.force)
}
