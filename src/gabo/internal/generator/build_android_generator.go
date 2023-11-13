package generator

import "fmt"

const _buildAndroidTask = `
      - name: Build with Gradle (dir %s)
        working-directory: "%s"
        run: |
          chmod +x gradlew
          ./gradlew buildDebug
`

func generateBuildAndroidYaml(repoDir string) (*string, error) {
	dirs, err := getDirsContaining(repoDir, "gradlew")
	if err != nil {
		return nil, err
	}
	str := _buildAndroidTemplate
	for _, dir := range dirs {
		str += fmt.Sprintf(_buildAndroidTask, dir, dir)
	}
	return &str, nil
}
