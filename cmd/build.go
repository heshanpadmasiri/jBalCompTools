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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Build project or file",
	Run: func(cmd *cobra.Command, args []string) {
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		} else {
			if viper.GetBool("file_comp") {
				fmt.Println("Please provide a file to compile")
				os.Exit(1)
			}
			targetPath = CurrentWorkingDir()
		}
		command, err := CreateCommand(viper.GetString("sourcePath"), viper.GetString("version"), targetPath, Build,
			viper.GetBool("remote_comp"))
		ConsumeError(err)
		if viper.GetBool("bench_comp") {
			result, err := BenchmarkCommand(&command)
			ConsumeError(err)
			PrettyPrintBenchmarkResult(result)
		} else {
			err = ExecuteCommand(&command)
			ConsumeError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolP("file", "f", false, "Run the given file")
	buildCmd.Flags().BoolP("remote", "r", false, "Remote debug the compiler")
	buildCmd.Flags().BoolP("benchmark", "b", false, "Benchmark the compiler")

	viper.BindPFlag("file_comp", buildCmd.Flags().Lookup("file"))
	viper.BindPFlag("remote_comp", buildCmd.Flags().Lookup("remote"))
	viper.BindPFlag("bench_comp", buildCmd.Flags().Lookup("benchmark"))
}
