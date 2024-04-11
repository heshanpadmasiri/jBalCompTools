package cmd

import (
	"testing"
)

func TestCreateCommand(t *testing.T) {
	sourcePath := "/path/to/source"
	version := "1.0.0"
	command := "build"
	targetPath := "/path/to/target"
	remoteDebug := true

	cmd, err := CreateCommand(sourcePath, version, command, targetPath, remoteDebug)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBalPath := BalPath(sourcePath, version)
	if cmd.Path != expectedBalPath {
		t.Errorf("Expected balPath to be %s, but got %s", expectedBalPath, cmd.Path)
	}

	expectedArgs := []string{expectedBalPath, command, targetPath, "--debug 5005"}
	if !stringSlicesEqual(cmd.Args, expectedArgs) {
		t.Errorf("Expected args to be %v, but got %v", expectedArgs, cmd.Args)
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
