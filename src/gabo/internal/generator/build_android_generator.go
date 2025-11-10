package generator

import (
	"fmt"
	"strings"
)

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
	var strSb19 strings.Builder
	for _, dir := range dirs {
		strSb19.WriteString(fmt.Sprintf(_buildAndroidTask, dir, dir))
	}
	str += strSb19.String()
	return &str, nil
}
