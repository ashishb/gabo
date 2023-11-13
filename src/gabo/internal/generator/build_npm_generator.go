package generator

import "fmt"

const _setupNodeJsTask = `
      - name: "Setup Node.js for %s"
        uses: actions/setup-node@v4
        with:
          node-version: "latest"
          cache: "npm"
          cache-dependency-path: "%s/package-lock.json"
`

const _buildNpmTask = `
      - name: "Build %s using NPM"
        working-directory: %s
        run: npm install && npm run build
`

func generateBuildNpmYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "package-lock.json")
	if err != nil {
		return nil, err
	}
	str := _buildNpmTemplate
	for _, dir := range dirs {
		str += fmt.Sprintf(_setupNodeJsTask, dir, dir)
		str += fmt.Sprintf(_buildNpmTask, dir, dir)
	}
	return &str, nil
}
