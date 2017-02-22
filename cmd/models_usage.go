package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stepanselyuk/gorm-play/models_usage"
)

// modelsUsageCmd represents the migrate command
var modelsUsageCmd = &cobra.Command{
	Use: "models-usage",
	Run: func(cmd *cobra.Command, args []string) {

		models_usage.FindCity()
	},
}

func init() {
	RootCmd.AddCommand(modelsUsageCmd)
	modelsUsageCmd.PersistentFlags().String("table", "example", "An example of a flag.")
}
