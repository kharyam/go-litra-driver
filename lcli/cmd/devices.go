package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List connected Litra devices",
	Long:  `Lists all connected Litra devices (Glow and Beam) with their indices.`,
	Run: func(cmd *cobra.Command, args []string) {
		devices := libImpl.ListDevices()
		if len(devices) == 0 {
			fmt.Println("No Litra devices found.")
			return
		}
		for _, d := range devices {
			fmt.Printf("  %d: Litra %s (serial: %s)\n", d.Index, d.Name, d.Serial)
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
