/*
Copyright © 2024 Heshan Padmasiri <hpheshan@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// disCmd represents the dis command
var disCmd = &cobra.Command{
	Use:   "dis <path>",
	Short: "Compile and dissemble a given file or project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the path to ballerina source/project to dissemble")
			os.Exit(1)
		}
		compileAndDissemble(args[0])
	},
}

func compileAndDissemble(path string) {
	compileTarget(viper.GetString("sourcePath"), viper.GetString("version"), path)
	jarName := getExpectedOutput(path)
	if _, err := os.Stat(jarName); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s not found\n", jarName)
		os.Exit(1)
	}
	createDisDir()
	moveJarToDisDir(jarName)
	disassemble(jarName)
}

func moveJarToDisDir(jarName string) {
	disPath := filepath.Join("dis", jarName)
	if err := os.Rename(jarName, disPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error moving jar file: %v\n", err)
		os.Exit(1)
	}
}

func createDisDir() {
	disDir := "dis"
	if _, err := os.Stat(disDir); err == nil {
		if err := os.RemoveAll(disDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting existing dis directory: %v\n", err)
			os.Exit(1)
		}
	}
	if err := os.Mkdir(disDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dis directory: %v\n", err)
		os.Exit(1)
	}
}

func disassemble(path string) {
	fmt.Println("Disassembling jar file...")
	cmd := exec.Command("jar", "-xf", path)
	cmd.Dir = "dis"
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running jar command: %v\n", err)
		os.Exit(1)
	}
}

func compileTarget(sourcePath, version, targetPath string) {
	command, err := CreateCommand(sourcePath, version, targetPath, Build, false)
	ConsumeError(err)
	err = ExecuteCommand(&command)
	ConsumeError(err)
}

func getExpectedOutput(path string) string {
	if isBallerinaProject(path) {
		return getProjectExpectedOutput(path)
	}
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, ".bal")
	fileName += ".jar"
	return fileName
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

func init() {
	rootCmd.AddCommand(disCmd)
}
