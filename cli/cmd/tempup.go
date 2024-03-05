/*
Copyright © 2024 Khary A. Mendez

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
	"strconv"

	"github.com/kharyam/go-litra-driver/lib"
	"github.com/spf13/cobra"
)

// tempupCmd represents the tempup command
var tempupCmd = &cobra.Command{
	Use:   "tempup",
	Short: "Increments the temperature by the amount specified",
	Long: `Increments the light temperature by the amount specified:

# Increment temperature by 100k
lcli tempup 100`,
	Run: func(cmd *cobra.Command, args []string) {
		temp, err := strconv.Atoi(args[0])
		if err != nil {
			temp = -1
		}
		if temp < 0 {
			fmt.Printf("Temperature increment must be a value greater than 0, not %s", args[0])
		} else {
			lib.LightTempUp(temp)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(tempupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tempupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tempupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
