package generator

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

const _validateSchemaTask = `
      - name: Validate OpenAPI schema in %s
        uses: mpetrunic/swagger-cli-action@v1.0.0
        with:
          command: 'validate %s/%s'
`

var _openAPIFileList = []string{"openapi.yaml", "openapi.yml", "openapi.json"}

func generateOpenAPISchemaValidator(repoDir string) (*string, error) {
	found := false
	template := _generateOpenAPISchemaValidatorTemplate
	for _, openAPIFile := range _openAPIFileList {
		result, err := generateOpenAPISchemaValidatorInternal(template, repoDir, openAPIFile)
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

func generateOpenAPISchemaValidatorInternal(template string, repoDir string, openAPIFile string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, openAPIFile)
	if err != nil {
		return nil, err
	}
	str := template
	for _, dir := range dirs {
		str += fmt.Sprintf(_validateSchemaTask, dir, dir, openAPIFile)
	}
	return &str, nil
}
