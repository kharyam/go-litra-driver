package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggles the light on or off",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, currentPower := libImpl.ReadCurrentState(deviceIndex)

		if currentPower == 1 {
			libImpl.LightOff(deviceIndex)
			fmt.Println("Light turned off")
		} else {
			libImpl.LightOn(deviceIndex)
			fmt.Println("Light turned on")
		}
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}
