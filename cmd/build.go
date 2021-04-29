package cmd

import (
	"github.com/lewis-od/lambda-build/pkg/command"
	"github.com/spf13/cobra"
)

var buildAndRun = command.NewBuildAndUploadCommand(orchestrator, tfExec, filesystem, printer)

var noBuild *bool
var buildCmd = &cobra.Command{
	Use:   "upload [version] [lambdas to build (optional)]",
	Short: "Build and upload lambdas",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buildAndRun.Run(args[0], args[1:], *noBuild)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	noBuild = buildCmd.Flags().Bool("no-build", false, "Don't build before uploading")
}
