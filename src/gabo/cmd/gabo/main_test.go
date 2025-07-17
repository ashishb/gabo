package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// A simple test that generates a docker linter for testing
func TestGenerateALinter(t *testing.T) {
	t.Parallel()
	dirPath := t.TempDir()
	t.Logf("Dir name is %s", dirPath)
	setFlagOrFail(t, "dir", dirPath)
	setFlagOrFail(t, "mode", _modeAnalyze)
	// Expected as this dir does not have ".git"
	require.Error(t, validateGitDir())
	// Now make it a git dir
	err := os.Mkdir(filepath.Join(dirPath, ".git"), 0o755)
	require.NoError(t, err)
	require.NoError(t, validateGitDir())

	// Now, do analyze again, verify it won't fail
	main()
	// Now, add a basic generate test for docker linting
	setFlagOrFail(t, "mode", _modeGenerate)
	setFlagOrFail(t, "for", "lint-docker")
	main()
	require.FileExists(t, filepath.Join(dirPath, ".github/workflows/lint-docker.yaml"))
}

func setFlagOrFail(t *testing.T, flagName string, value string) {
	t.Helper()
	err := flag.Set(flagName, value)
	require.NoError(t, err)
}
