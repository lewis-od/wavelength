package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all lambdas",
	Run: func(cmd *cobra.Command, args []string) {
		lambdas, err := finder.FindLambdas([]string{})
		if err != nil {
			printer.PrintErr(err)
			os.Exit(1)
		}

		if len(lambdas) == 0 {
			printer.Println("🕵️ No lambdas found")
			return
		}

		printer.Println("🕵️ Found lambdas:")
		for _, lambda := range lambdas {
			printer.Println(lambda)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
