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

func CreateCommand(sourcePath, version, command, targetPath string, remoteDebug bool) (exec.Cmd, error) {
	balPath := BalPath(sourcePath, version)
	// if !fileExists(balPath) {
	// 	return exec.Cmd{}, fmt.Errorf("bal executable not found at %s try running 'jBalCompTools build'", balPath)
	// }
	// if !fileExists(targetPath) {
	// 	return exec.Cmd{}, fmt.Errorf("target path not found at %s", targetPath)
	// }
	args := []string{ command, targetPath }
	if remoteDebug {
		args = append(args, "--debug 5005")
	}
	cmd := exec.Command(balPath, args...)
	return *cmd, nil
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

// TODO: remove this from run
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
