package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stepanselyuk/gorm-play/db_mappings"
)

// initializeCmd represents the migrate command
var xmlLoadCmd = &cobra.Command{
	Use: "xml_load",
	Run: func(cmd *cobra.Command, args []string) {

		db_mappings.GetProductMapping()
	},
}

func init() {
	RootCmd.AddCommand(xmlLoadCmd)
}
