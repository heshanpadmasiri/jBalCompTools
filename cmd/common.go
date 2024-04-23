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
	"time"
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
		// create copy of cmd
		cmdCopy := *cmd
		err := ExecuteCommand(&cmdCopy)
		elapsed := time.Since(start)
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
