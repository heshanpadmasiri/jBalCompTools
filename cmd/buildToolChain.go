// Copyright (c) 2024 Heshan Padmasiri
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildToolsCmd = &cobra.Command{
	Use:   "buildTools",
	Short: "Build jBallerina compiler",
	Run: func(cmd *cobra.Command, args []string) {
		err := BuildCompiler(viper.GetString("sourcePath"), viper.GetString("flags"))
		ConsumeError(err)
	},
}

func init() {
	rootCmd.AddCommand(buildToolsCmd)
	buildToolsCmd.Flags().String("flags", "buildTools -x check", "Flags to pass to the gradle wrapper")
	viper.BindPFlag("flags", buildToolsCmd.Flags().Lookup("flags"))
}
