package cmd

import (
	"fmt"

	"github.com/kharyam/go-litra-driver/config"
	"github.com/kharyam/go-litra-driver/lib"
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggles the light on or off",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, currentPower := config.ReadCurrentState()

		if currentPower == 1 {
			lib.LightOff()
			fmt.Println("Light turned off")
		} else {
			lib.LightOn()
			fmt.Println("Light turned on")
		}
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}
