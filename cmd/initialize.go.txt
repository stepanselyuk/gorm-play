package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stepanselyuk/gorm-play/initialize"
)

// initializeCmd represents the migrate command
var initializeCmd = &cobra.Command{
	Use: "initialize",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("initialize called")
		initialize.Example()
	},
}

func init() {
	RootCmd.AddCommand(initializeCmd)
	initializeCmd.PersistentFlags().String("table", "example", "An example of a flag.")
}
