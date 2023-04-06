package generator

import "fmt"

const _buildDockerTask = `
      - name: Docker build using %s/Dockerfile
        run: |
          cd "%s"
          DOCKER_BUILDKIT=1 docker build -f Dockerfile .
`

func generateBuildDockerYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "Dockerfile")
	if err != nil {
		return nil, err
	}
	str := _buildDockerTemplate
	for _, dir := range dirs {
		str += fmt.Sprintf(_buildDockerTask, dir, dir)
	}
	return &str, nil
}
