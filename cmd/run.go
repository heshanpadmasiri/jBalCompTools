// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run the current project or file",
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		} else {
			if (viper.GetBool("file")) {
				fmt.Println("Please provide a file to run")
				os.Exit(1)
			}
			targetPath = CurrentWorkingDir()
		}
		command, err := CreateCommand(viper.GetString("sourcePath"), viper.GetString("version"), "run", targetPath,
			viper.GetBool("remote"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ExecuteCommand(command)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("file", "f", false, "Run the given file")
	runCmd.Flags().BoolP("remote", "r", false, "Remote debug the runtime")
	viper.BindPFlag("file", runCmd.Flags().Lookup("file"))
	viper.BindPFlag("remote", runCmd.Flags().Lookup("remote"))
}
