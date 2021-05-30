package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const version string = "v1.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wavelength",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
