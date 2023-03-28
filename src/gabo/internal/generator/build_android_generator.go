package generator

import "fmt"

const _buildAndroidTask = `
      - name: Build with Gradle (dir %s)
        run: |
          cd "%s"
          chmod +x gradlew
          ./gradlew buildDebug
`

func generateBuildAndroidYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "gradlew")
	if err != nil {
		return nil, err
	}
	str := _buildAndroidYaml
	for _, dir := range dirs {
		str += fmt.Sprintf(_buildAndroidTask, dir, dir)
	}
	return &str, nil
}
