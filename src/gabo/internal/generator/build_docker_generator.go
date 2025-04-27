package generator

import "fmt"

const _buildDockerTask = `
      - name: Docker build using %s/Dockerfile
        working-directory: "%s"
        run: |
          DOCKER_BUILDKIT=1 docker buildx build --cache-from type=gha --cache-to type=gha -t tmp1 -f Dockerfile .
`

const _checkDockerWithDiveTask = `
	  - name: Check Docker image (%s/Dockerfile) with dive for wasted space
        working-directory: "%s"
		run: |
			docker run --rm \
				-v /var/run/docker.sock:/var/run/docker.sock \
				--env=CI=true \
				--network none \
				docker.io/wagoodman/dive:latest \
				tmp1
`

func generateBuildDockerYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "Dockerfile")
	if err != nil {
		return nil, err
	}
	str := _buildDockerTemplate
	for _, dir := range dirs {
		str += fmt.Sprintf(_buildDockerTask+"\n\n", dir, dir)
		str += fmt.Sprintf(_checkDockerWithDiveTask, dir, dir)
	}
	return &str, nil
}
