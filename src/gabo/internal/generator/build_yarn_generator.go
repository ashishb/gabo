package generator

import (
	"fmt"
	"strings"
)

const _setupYarnTask = `
      - name: "Setup Node.js for %s"
        uses: actions/setup-node@v4
        with:
          node-version: "latest"
          cache: "yarn"
          cache-dependency-path: "%s/yarn.lock"
`

const _buildYarnTask = `
      - name: "Build %s using Yarn"
        working-directory: %s
        run: yarn && yarn build
`

func generateBuildYarnYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "yarn.lock")
	if err != nil {
		return nil, err
	}
	str := _buildYarnTemplate
	var strSb26 strings.Builder
	for _, dir := range dirs {
		strSb26.WriteString(fmt.Sprintf(_setupYarnTask, dir, dir))
		strSb26.WriteString(fmt.Sprintf(_buildYarnTask, dir, dir))
	}
	str += strSb26.String()
	return &str, nil
}
