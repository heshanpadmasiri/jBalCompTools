// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Command string

const (
	Run   Command = "run"
	Build Command = "build"
)

func CreateCommand(sourcePath, version, targetPath string, command Command, remoteDebug bool) (exec.Cmd, error) {
	balPath := BalPath(sourcePath, version)
	switch command {
	case Run:
		return createRunCommand(balPath, targetPath, remoteDebug), nil
	case Build:
		return createBuildCommand(balPath, targetPath, remoteDebug), nil
	default:
		return exec.Cmd{}, fmt.Errorf("unknown command: %s", command)
	}
}

func createBuildCommand(balPath, targetPath string, remoteDebug bool) exec.Cmd {
	args := []string {"build", targetPath}
	cmd := exec.Command(balPath, args...)
	if remoteDebug {
		cmd.Env = append(os.Environ(), "BAL_JAVA_DEBUG=5005")
	}
	return *cmd
}

func createRunCommand(balPath, targetPath string, remoteDebug bool) exec.Cmd {
	args := []string{"run"}
	if remoteDebug {
		args = append(args, "--debug", "5005")
	}
	args = append(args, targetPath)
	return *exec.Command(balPath, args...)
}

func ExecuteCommand(cmd exec.Cmd) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func BalPath(srcPath, version string) string {
	return filepath.Join(srcPath, "distribution", "zip", "jballerina-tools", "build", "extracted-distributions",
		fmt.Sprintf("jballerina-tools-%s", version), "bin", "bal")
}

func CurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		os.Exit(1)
	}
	return dir
}
