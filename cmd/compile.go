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

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [path]",
	Short: "Compile project or file",
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
		command, err := CreateCommand(viper.GetString("sourcePath"), viper.GetString("version"), targetPath, Build,
			viper.GetBool("remote"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ExecuteCommand(command)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.Flags().BoolP("file", "f", false, "Run the given file")
	compileCmd.Flags().BoolP("remote", "r", false, "Remote debug the runtime")

	viper.BindPFlag("file", compileCmd.Flags().Lookup("file"))
	viper.BindPFlag("remote", compileCmd.Flags().Lookup("remote"))
}
