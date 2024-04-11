package cmd

import (
	"testing"
)

func TestCreateCommandRun(t *testing.T) {
	sourcePath := "/path/to/source"
	version := "1.0.0"
	command := Run
	targetPath := "/path/to/target"
	remoteDebug := true

	cmd, err := CreateCommand(sourcePath, version, targetPath, command, remoteDebug)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBalPath := BalPath(sourcePath, version)

	expectedArgs := []string{expectedBalPath, string(command),  "--debug", "5005", targetPath}
	if !stringSlicesEqual(cmd.Args, expectedArgs) {
		t.Errorf("Expected args to be %v, but got %v", expectedArgs, cmd.Args)
	}
}

func TestCreateCommandBuild(t *testing.T) {
	sourcePath := "/path/to/source"
	version := "1.0.0"
	command := Build
	targetPath := "/path/to/target"
	remoteDebug := true

	cmd, err := CreateCommand(sourcePath, version, targetPath, command, remoteDebug)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBalPath := BalPath(sourcePath, version)
	expectedArgs := []string{expectedBalPath, string( command), targetPath}
	if !stringSlicesEqual(cmd.Args, expectedArgs) {
		t.Errorf("Expected args to be %v, but got %v", expectedArgs, cmd.Args)
	}

	if !stringArrayContains(cmd.Env, "BAL_JAVA_DEBUG=5005") {
		t.Errorf("Expected env to contain BAL_JAVA_DEBUG=5005, but got %v", cmd.Env)
	}
}

func stringArrayContains(arr []string, s string) bool {
	for _, e := range arr {
		if e == s {
			return true
		}
	}
	return false
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
