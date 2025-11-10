package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const _validateSchemaTask = `
      # Ref: https://github.com/daveshanley/vacuum
      - name: Validate OpenAPI schema
        working-directory: head
        run:
          docker run --rm -v $PWD:/work:ro dshanley/vacuum lint -d %s

      # Ref: https://github.com/Tufin/oasdiff
      - name: Running OpenAPI Spec diff action
        env:
          BASE: ./base/%s
          REVISION: ./head/%s
        run: |
          docker run --rm -v $PWD:/code:ro -t tufin/oasdiff breaking --fail-on WARN -f githubactions "${BASE}" "${REVISION}"
`

var _openAPIFileList = []string{"openapi.yaml", "openapi.yml", "**/openapi.json"}

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

func generateOpenAPISchemaValidatorInternal(
	template string, repoDir string, openAPIFile string,
) (*string, error) {
	dirs, err := getDirsContaining(repoDir, openAPIFile)
	if err != nil {
		return nil, err
	}
	str := template
	var strSb57 strings.Builder
	for _, dir := range dirs {
		path := fmt.Sprintf("%s/%s", dir, openAPIFile)
		strSb57.WriteString(fmt.Sprintf(_validateSchemaTask, path, path, path))
	}
	str += strSb57.String()
	return &str, nil
}
