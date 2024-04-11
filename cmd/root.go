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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jBalCompTools",
	Short: "Collection of useful commands for jBallerina compiler debugging",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/jBalCompTools")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
	}

	rootCmd.PersistentFlags().StringP("sourcePath", "s", viper.GetString("defaultSourcePath"), "Path to jBallerina source code")
	viper.BindPFlag("sourcePath", rootCmd.PersistentFlags().Lookup("sourcePath"))

	rootCmd.PersistentFlags().StringP("version", "v", viper.GetString("defaultVersion"), "Version of jBallerina")
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
}
