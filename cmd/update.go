package cmd

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/service"

	"github.com/spf13/cobra"
)

var updateService = service.NewUpdateService(finder, updater, printer, &projectName)

var roleId = ""

var updateCmd = &cobra.Command{
	Use:   "update [version] [lambdas to build]",
	Short: "Updates the code used for the specified lambdas",
	Long: `Updates the specified lambda's code with the artifact at <version>/<lambda name>.zip in S3.

If no lambdas are specified, all will be updated.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		roleToAssume := &builder.Role{RoleID: roleId}
		if roleId == "" {
			roleToAssume = nil
		}
		updateService.Run(args[0], args[1:], roleToAssume)
	},
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&roleId, "assume-role", "a", "", "AWS role to assume")
	rootCmd.AddCommand(updateCmd)
}
