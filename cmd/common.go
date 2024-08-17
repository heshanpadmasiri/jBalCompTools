// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"log"
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
	Test  Command = "test"
)

func CreateJarRunCommand(jarPath string) exec.Cmd {
	return *exec.Command("java", "-jar", jarPath)
}

func CreateCommand(sourcePath, version, targetPath string, command Command, remoteDebug bool, args ...string) (exec.Cmd, error) {
	ConsumeError(buildCompilerIfneeded(sourcePath, version))
	balPath := BalPath(sourcePath, version)
	switch command {
	case Run:
		return createRunCommand(balPath, targetPath, remoteDebug, args...), nil
	case Build:
		return createBuildCommand(balPath, targetPath, remoteDebug, args...), nil
	case Test:
		return createTestCommand(balPath, targetPath, remoteDebug, args...), nil
	default:
		return exec.Cmd{}, fmt.Errorf("unknown command: %s", command)
	}
}

func createBuildCommand(balPath, targetPath string, remoteDebug bool, extraArgs ...string) exec.Cmd {
	args := []string{"build"}
	args = append(args, extraArgs...)
	args = append(args, targetPath)
	cmd := exec.Command(balPath, args...)
	if remoteDebug {
		cmd.Env = append(os.Environ(), "BAL_JAVA_DEBUG=5005")
	}
	return *cmd
}

func createTestCommand(balPath, targetPath string, remoteDebug bool, extraArgs ...string) exec.Cmd {
	return createExecCommand(balPath, targetPath, "test", remoteDebug)
}

func createRunCommand(balPath, targetPath string, remoteDebug bool, extraArgs ...string) exec.Cmd {
	return createExecCommand(balPath, targetPath, "run", remoteDebug)
}

func createExecCommand(balPath, targetPath, runCommand string, remoteDebug bool) exec.Cmd {
	args := []string{runCommand}
	if remoteDebug {
		args = append(args, "--debug", "5005")
	}
	args = append(args, targetPath)
	return *exec.Command(balPath, args...)
}

func ExecuteCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Executing cmd")
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

func buildCompilerIfneeded(sourcePath, version string) error {
	balPath := BalPath(sourcePath, version)
	if !compilerExists(balPath) {
		return BuildCompiler(sourcePath, "build -x check")
	}
	return nil;
}

func BuildCompiler(path, flags string) error {
	args := strings.Split(strings.Trim(flags, " "), " ")
	cmd := exec.Command("./gradlew", args...)
	cmd.Dir = path
	return ExecuteCommand(cmd)
}

func compilerExists(balPath string) bool {
	_, err := os.Stat(balPath)
	return err == nil;
}
