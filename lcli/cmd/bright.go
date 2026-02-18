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

	"github.com/spf13/cobra"
)

// brightCmd represents the bright command
var brightCmd = &cobra.Command{
	Use:   "bright",
	Short: "Sets the brightness level (0-100)",
	Long:  `Sets the brightness level of all lights. Specify a value level between 0 and 100`,
	Run: func(cmd *cobra.Command, args []string) {
		bright, err := strconv.Atoi(args[0])
		if err != nil {
			bright = -1
		}
		if bright < 0 || bright > 100 {
			fmt.Printf("Brightness must be a value between 0 and 100, not %s", args[0])
		} else {
			libImpl.LightBrightness(deviceIndex, bright)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(brightCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brightCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brightCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
