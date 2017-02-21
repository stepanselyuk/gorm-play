package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stepanselyuk/gorm-play/migrate"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("migrate called")
		migrate.Example()
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().String("table", "example", "An example of a flag.")
}
