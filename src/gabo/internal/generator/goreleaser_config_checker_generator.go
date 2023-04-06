package generator

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

const _goReleaserCheckerTask = `
      - name: Check Go Releaser config is valid
        run: goreleaser check --config %s/%s

      - name: Build (not release) binaries with Go Releaser
        run: goreleaser build --snapshot --clean
`

var _goReleaserConfigFiles = []string{
	"goreleaser.yaml", "goreleaser.yml", ".goreleaser.yaml", ".goreleaser.yml",
}

func generateGoReleaserConfigCheckerYaml(repoDir string) (*string, error) {
	found := false
	template := _checkGoReleaserConfigTemplate
	for _, openAPIFile := range _goReleaserConfigFiles {
		result, err := generateGoReleaserCheckerInternal(template, repoDir, openAPIFile)
		if err != nil && !errors.Is(err, errNoSuchDir) {
			return nil, err
		}
		if errors.Is(err, errNoSuchDir) {
			log.Debug().Msgf("No dir containing %s found", openAPIFile)
			continue
		}
		template = *result
		found = true
	}
	if !found {
		return nil, errNoSuchDir
	}
	return &template, nil
}

func generateGoReleaserCheckerInternal(template, repoDir, releaserFile string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, releaserFile)
	if err != nil {
		return nil, err
	}
	str := template
	for _, dir := range dirs {
		str += fmt.Sprintf(_goReleaserCheckerTask, dir, releaserFile)
	}
	return &str, nil
}
