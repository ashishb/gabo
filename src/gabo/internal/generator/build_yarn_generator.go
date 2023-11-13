package generator

import "fmt"

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
	for _, dir := range dirs {
		str += fmt.Sprintf(_setupYarnTask, dir, dir)
		str += fmt.Sprintf(_buildYarnTask, dir, dir)
	}
	return &str, nil
}
