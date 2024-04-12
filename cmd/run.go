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
	Short: "Run project or file",
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		} else {
			if (viper.GetBool("file_run")) {
				fmt.Println("Please provide a file to run")
				os.Exit(1)
			}
			targetPath = CurrentWorkingDir()
		}
		command, err := CreateCommand(viper.GetString("sourcePath"), viper.GetString("version"), targetPath, Run,
			viper.GetBool("remote_run"))
		ConsumeError(err)
		err = ExecuteCommand(&command)
		ConsumeError(err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("file", "f", false, "Run the given file")
	runCmd.Flags().BoolP("remote", "r", false, "Remote debug the runtime")
	viper.BindPFlag("file_run", runCmd.Flags().Lookup("file"))
	viper.BindPFlag("remote_run", runCmd.Flags().Lookup("remote"))
}
