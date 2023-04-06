package generator

import (
	"fmt"
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

func (g Generator) Generate(optionStr string) error {
	if g.force {
		log.Warn().Msgf("Force overwrite is on, existing files will be over-written")
	}
	option, err := getOptionForFlag(optionStr)
	if err != nil {
		return err
	}
	str, err := (*option).GetYamlConfig(g.dir)
	if err != nil {
		return err
	}
	return writeOrWarn((*option).GetOutputFileName(g.dir), *str, g.force)
}

func getOptionForFlag(target string) (*Option, error) {
	for _, option := range GetOptions() {
		if option.FlagName() == target {
			return &option, nil
		}
	}
	return nil, fmt.Errorf("invalid option: %s", target)
}
