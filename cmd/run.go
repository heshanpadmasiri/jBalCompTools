// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run the current project or file",
	Run: func(cmd *cobra.Command, args []string) {
		balPath := BalPath(viper.GetString("sourcePath"), viper.GetString("version"))
		if !fileExists(balPath) {
			fmt.Printf("bal executable not found at %s try running 'jBalCompTools build'\n", balPath)
			os.Exit(1)
		}
		cwd := CurrentWorkingDir()
		isFile := viper.GetBool("file")
		isRemoteDebug := viper.GetBool("remote")
		if isFile {
			if len(args) == 0 {
				fmt.Println("When running a file path is required")
				os.Exit(1)
			}
			runFile(balPath, args[0], isRemoteDebug)
			return
		}
		targetPath := cwd
		if len(args) > 0 {
			targetPath = args[0]
		}
		runProject(balPath, targetPath, isRemoteDebug)
	},
}

func runFile(balPath, filePath string, remoteDebug bool) {
	args := []string{ "run", filePath }
	if remoteDebug {
		args = append(args, "--debug 5005")
	}
	cmd := exec.Command(balPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func runProject(balPath, projectPath string, remoteDebug bool) {
	args := []string{ "run"  }
	if remoteDebug {
		args = append(args, "--debug 5005")
	}
	cmd := exec.Command(balPath, args...)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("file", "f", false, "Run the given file")
	runCmd.Flags().BoolP("remote", "r", false, "Remote debug the runtime")
	viper.BindPFlag("file", runCmd.Flags().Lookup("file"))
	viper.BindPFlag("remote", runCmd.Flags().Lookup("remote"))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
