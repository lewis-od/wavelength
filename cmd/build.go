package cmd

import (
	"github.com/lewis-od/wavelength/internal/service"
	"github.com/spf13/cobra"
)

var buildAndRun = service.NewBuildAndUploadService(orchestrator, finder, printer)

var noBuild *bool
var buildCmd = &cobra.Command{
	Use:   "upload [version] [lambdas to build]",
	Short: "Builds and uploads lambdas",
	Long: `Builds the specified lambdas using lerna, then uploads the build artifact to S3 with the key <version>/<lambda name>.zip

If no lambdas are specified, all will be built and uploaded.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buildAndRun.Run(args[0], args[1:], *noBuild)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	noBuild = buildCmd.Flags().Bool("no-build", false, "Don't build before uploading")
}
