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
	"strings"
	"time"

	"github.com/BurntSushi/toml"
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

type BenchmarkResult struct {
	Iterations int
	AvgTime    time.Duration
	MinTime    time.Duration
	MaxTime    time.Duration
}

// TODO: also accept a configuration
func BenchmarkCommand(cmd *exec.Cmd) (BenchmarkResult, error) {
	nIterations := 10
	elapsedTimes := make([]time.Duration, nIterations)
	for i := 0; i < nIterations; i++ {
		start := time.Now()
		cmdCopy := *cmd
		cmdCopy.Stdout = nil
		cmdCopy.Stderr = os.Stderr
		err := cmdCopy.Run()
		elapsed := time.Since(start)
		fmt.Println("Iteration: ", i, "elapsed: ", elapsed)
		if err != nil {
			return BenchmarkResult{}, err
		}
		elapsedTimes[i] = elapsed
	}
	return analyzeElapsedTimes(elapsedTimes), nil
}

func analyzeElapsedTimes(elapsedTimes []time.Duration) BenchmarkResult {
	var sum time.Duration
	min := elapsedTimes[0]
	max := elapsedTimes[0]
	for _, elapsed := range elapsedTimes {
		sum += elapsed
		if elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
	}
	n := len(elapsedTimes)
	avg := sum / time.Duration(n)
	return BenchmarkResult{
		Iterations: n,
		AvgTime:    avg,
		MinTime:    min,
		MaxTime:    max,
	}
}

func PrettyPrintBenchmarkResult(result BenchmarkResult) {
	fmt.Printf("Ran %d iterations\n", result.Iterations)
	fmt.Printf("Average time: %v\n", result.AvgTime)
	fmt.Printf("Minimum time: %v\n", result.MinTime)
	fmt.Printf("Maximum time: %v\n", result.MaxTime)
}

func CompileTarget(sourcePath, version, targetPath string) {
	command, err := CreateCommand(sourcePath, version, targetPath, Build, false)
	ConsumeError(err)
	err = ExecuteCommand(&command)
	ConsumeError(err)
}

func GetExpectedOutput(path string) string {
	if isBallerinaProject(path) {
		return getProjectExpectedOutput(path)
	}
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, ".bal")
	fileName += ".jar"
	return fileName
}

func isBallerinaProject(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Fprintf(os.Stderr, "Error getting file info at path %s: %v\n", path, err)
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		return false
	}

	tomlPath := filepath.Join(path, "Ballerina.toml")
	_, err = os.Stat(tomlPath)
	return err == nil
}

func getProjectExpectedOutput(path string) string {
	balTomlPath := filepath.Join(path, "Ballerina.toml")
	tomlContent, err := os.ReadFile(balTomlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading Ballerina.toml file: %v\n", err)
		os.Exit(1)
	}

	var config struct {
		Package struct {
			Name string `toml:"name"`
		} `toml:"package"`
	}

	if err := toml.Unmarshal(tomlContent, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling Ballerina.toml file: %v\n", err)
		os.Exit(1)
	}

	name := config.Package.Name
	return name + ".jar"
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
