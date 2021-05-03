package cmd

import (
	"github.com/lewis-od/wavelength/internal/service"

	"github.com/spf13/cobra"
)

var updateService = service.NewUpdateService(finder, updater, printer, projectName)

var updateCmd = &cobra.Command{
	Use:   "update [version] [lambdas to build (optional)]",
	Short: "Update the code used for the specified lambdas",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		updateService.Run(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
