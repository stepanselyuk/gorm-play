package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stepanselyuk/gorm-play/db_mappings"
)

// initializeCmd represents the migrate command
var generateModelsCmd = &cobra.Command{
	Use: "generate-models",
	Run: func(cmd *cobra.Command, args []string) {

		db_mappings.GenerateGormModels()
	},
}

func init() {
	RootCmd.AddCommand(generateModelsCmd)
}
