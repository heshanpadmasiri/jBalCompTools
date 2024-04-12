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
	args := []string{"build", targetPath}
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

func ExecuteCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BalPath(srcPath, version string) string {
	return filepath.Join(srcPath, "distribution", "zip", "jballerina-tools", "build", "extracted-distributions",
		fmt.Sprintf("jballerina-tools-%s", version), "bin", "bal")
}

func CurrentWorkingDir() string {
	dir, err := os.Getwd()
	ConsumeError(err)
	return dir
}

func ConsumeError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
