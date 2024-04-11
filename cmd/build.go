// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build jBallerina compiler",
	Run: func(cmd *cobra.Command, args []string) {
		err := buildCompiler(viper.GetString("sourcePath"), viper.GetString("flags"))
		if (err != nil) {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().String("flags", "build -x check", "Flags to pass to the gradle wrapper")
	viper.BindPFlag("flags", buildCmd.Flags().Lookup("flags"))
}

func buildCompiler(path, flags string) error {
	// fmt.Println("DEBUG: ", path, flags)
	args := strings.Split(strings.Trim(flags, " "), " ")
	cmd := exec.Command("./gradlew", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("failed to execute build command: %v", err)
	}
	return nil
}
