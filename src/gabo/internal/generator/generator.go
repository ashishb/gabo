package generator

import (
	"github.com/rs/zerolog/log"
)

func Generate(option Option, force bool) error {
	if force {
		log.Warn().Msgf("Force overwrite is on, existing files will be over-written")
	}
	switch option {
	case _LintDocker:
		return writeOrWarn(option.getOutputFileName(), option.getYamlConfig(), force)
	}
	return nil
}
